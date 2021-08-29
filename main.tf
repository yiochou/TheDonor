variable "IMAGE_TAG" {
  type = string
}

resource "google_cloud_run_service" "default" {
  provider = google-beta

  name     = "parseandtweet"
  location = "asia-east1"

  template {
    spec {
      containers {
        image = "gcr.io/thedonor/main:${var.IMAGE_TAG}"
        env {
          name = "MONGODB_URI"
          value_from {
            secret_key_ref {
              name = "MONGODB_URI"
              key  = "1"
            }
          }
        }
        env {
          name = "MONGODB_USERNAME"
          value_from {
            secret_key_ref {
              name = "MONGODB_USERNAME"
              key  = "1"
            }
          }
        }
        env {
          name = "MONGODB_PASSWORD"
          value_from {
            secret_key_ref {
              name = "MONGODB_PASSWORD"
              key  = "1"
            }
          }
        }
        env {
          name = "TWITTER_CONSUMER_KEY"
          value_from {
            secret_key_ref {
              name = "TWITTER_CONSUMER_KEY"
              key  = "1"
            }
          }
        }
        env {
          name = "TWITTER_CONSUMER_SECRET"
          value_from {
            secret_key_ref {
              name = "TWITTER_CONSUMER_SECRET"
              key  = "1"
            }
          }
        }
        env {
          name = "TWITTER_TOKEN"
          value_from {
            secret_key_ref {
              name = "TWITTER_TOKEN"
              key  = "1"
            }
          }
        }
        env {
          name = "TWITTER_TOKEN_SECRET"
          value_from {
            secret_key_ref {
              name = "TWITTER_TOKEN_SECRET"
              key  = "1"
            }
          }
        }
      }
    }
  }

  metadata {
    annotations = {
      generated-by                      = "magic-modules"
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      metadata.0.annotations,
    ]
  }
}

resource "google_cloud_scheduler_job" "job" {
  name             = "tweet-cases"
  description      = "tweet parsed cases"
  schedule         = "0 10 * * *"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://parseandtweet-vik6b4pw7a-de.a.run.app/"

    oidc_token {
      service_account_email = "scheduler@thedonor.iam.gserviceaccount.com"
    }
  }
}

provider "google" {
  project = "thedonor"
  region  = "asia-east1"
}

provider "google-beta" {
  project = "thedonor"
  region  = "asia-east1"
}
