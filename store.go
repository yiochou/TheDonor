package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db     *mongo.Database
	logger log.Logger
}

func ConnectMongoDB(config Config) (*mongo.Database, error) {
	credential := options.Credential{
		Username: config.MongodbUsername,
		Password: config.MongodbPassword,
	}

	clientOptions := options.Client().ApplyURI(config.MongodbUri).SetAuth(credential)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	log.Info("connected to MongoDB!")

	db := client.Database(config.Database)

	return db, nil
}

func NewStore(db *mongo.Database, logger log.Logger) Store {
	return Store{
		db:     db,
		logger: logger,
	}
}

func (store *Store) CaseCollection() *mongo.Collection {
	return store.db.Collection("cases")
}

func (store *Store) InsertCases(cases []*Case) error {
	toSave := make([]interface{}, len(cases))
	for i, c := range cases {
		toSave[i] = c
	}

	insertResult, err := store.CaseCollection().InsertMany(context.TODO(), toSave)
	if err != nil {
		return err
	}

	store.logger.Info("cases inserted, _ids: ", insertResult)

	return nil
}

func (store *Store) FilterNewExistCases(cases []*Case) ([]*Case, error) {
	options := options.Find().SetProjection(bson.M{"id": 1})
	filter := bson.M{"id": bson.M{"$in": filterCaseIds(cases)}}

	cur, err := store.CaseCollection().Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	existedCaseMap := map[string]bool{}

	for cur.Next(context.TODO()) {
		var c Case
		err := cur.Decode(&c)
		if err != nil {
			store.logger.Error(err)
			return nil, err
		}

		existedCaseMap[c.Id] = true
	}

	if err := cur.Err(); err != nil {
		store.logger.Error(err)
		return nil, err
	}

	cur.Close(context.TODO())

	var newCases []*Case
	for _, c := range cases {
		if _, ok := existedCaseMap[c.Id]; !ok {
			newCases = append(newCases, c)
		}
	}

	return newCases, nil

}
func filterCaseIds(cases []*Case) []string {
	var ids []string

	for _, c := range cases {
		ids = append(ids, c.Id)
	}

	return ids
}

func (store *Store) InsertCasesIfNotExists(cases []*Case) ([]*Case, error) {
	newCases, err := store.FilterNewExistCases(cases)
	if err != nil {
		return nil, err
	}
	if newCases == nil {
		return nil, nil
	}

	err = store.InsertCases(newCases)

	if err != nil {
		return nil, err
	}

	return newCases, nil
}
