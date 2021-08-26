package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	mongoClient = ConnectMongoDB()
	Ping()
}

func ConnectMongoDB() *mongo.Client {
	credential := options.Credential{
		Username: viper.GetString("MONGODB_USERNAME"),
		Password: viper.GetString("MONGODB_PASSWORD"),
	}
	uri := viper.GetString("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getMongoClient() *mongo.Client {
	if mongoClient == nil {
		ConnectMongoDB()
	}
	return mongoClient
}

func getCaseCollection() *mongo.Collection {
	return getMongoClient().Database("the_donor").Collection("cases")
}

func InsertCases(cases []*Case) error {
	toSave := make([]interface{}, len(cases))
	for i, c := range cases {
		toSave[i] = c
	}
	insertResult, err := getCaseCollection().InsertMany(context.TODO(), toSave)

	if err != nil {
		return err
	}

	log.Info("cases inserted, _ids: ", insertResult)

	return nil
}

func FilterNewExistCases(cases []*Case) ([]*Case, error) {
	options := options.Find().SetProjection(bson.M{"id": 1})
	filter := bson.M{"id": bson.M{"$in": filterCaseIds(cases)}}

	cur, err := getCaseCollection().Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	existedCaseMap := map[string]bool{}

	for cur.Next(context.TODO()) {
		var c Case
		err := cur.Decode(&c)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		existedCaseMap[c.Id] = true
	}

	if err := cur.Err(); err != nil {
		log.Error(err)
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

func Ping() {
	err := getMongoClient().Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func InsertCasesIfNotExists(cases []*Case) ([]*Case, error) {
	newCases, err := FilterNewExistCases(cases)
	if err != nil {
		return nil, err
	}
	if newCases == nil {
		return nil, nil
	}

	err = InsertCases(newCases)

	if err != nil {
		return nil, err
	}

	return newCases, nil
}
