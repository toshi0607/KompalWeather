// https://www.terraform.io/docs/providers/google/r/cloud_run_service.html
resource "google_cloud_run_service" "kompal-weather-dev" {
  name     = "kompal-weather-dev2"
  location = var.gcp_region

  autogenerate_revision_name = false

  traffic {
    latest_revision = true
    percent         = 100
  }

  template {
    spec {
      container_concurrency = 10
      service_account_name  = google_service_account.kompal_weather.email
      timeout_seconds       = 300

      containers {
        image = "asia.gcr.io/terraform-toshi0607/github_toshi0607_kompalweather"
        env {
          name  = "SPREAD_SHEET_ID"
          value = "1f_O-yUZ3KFFlduHMsQCCZUd4yVNS82zE0lePaDe1WoE"
        }
        env {
          name  = "SHEET_ID"
          value = "0"
        }
        env {
          name  = "GCP_PROJECT_ID"
          value = var.gcp_project
        }
        env {
          name  = "SERVER_PORT"
          value = "8080"
        }
        env {
          name  = "SLACK_CHANNEL_NAMES"
          value = "dev"
        }
        env {
          name  = "SLACK_WEBHOOK_URL_SECRET_NAME"
          value = google_secret_manager_secret.web_hook_url.secret_id
        }
        env {
          name  = "KOMPAL_URL_SECRET_NAME"
          value = google_secret_manager_secret.kompal_url.secret_id
        }
        env {
          name  = "TWITTER_ACCESS_TOKEN_SECRET_NAME"
          value = google_secret_manager_secret.twitter_access_token.secret_id
        }
        env {
          name  = "TWITTER_ACCESS_TOKEN_SECRET_SECRET_NAME"
          value = google_secret_manager_secret.twitter_access_token_secret.secret_id
        }
        env {
          name  = "TWITTER_API_KEY_SECRET_NAME"
          value = google_secret_manager_secret.twitter_api_key.secret_id
        }
        env {
          name  = "TWITTER_API_KEY_SECRET_SECRET_NAME"
          value = google_secret_manager_secret.twitter_api_key_secret.secret_id
        }
        env {
          name  = "VERSION"
          value = "0.1.0"
        }
        env {
          name  = "SERVICE_NAME"
          value = "KompalWeather"
        }
        env {
          name  = "ENVIRONMENT"
          value = "development"
        }
        resources {
          limits = {
            cpu    = "1000m"
            memory = "256Mi"
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template.0.spec.0.containers.0.image]
  }

  depends_on = [
    google_project_iam_member.kompal_weather_is_logging_log_writer,
    google_project_iam_member.kompal_weather_is_secret_manager_admin
  ]
}

resource "google_cloud_run_service" "kompal_weather_visualizer_dev" {
  name     = "kompal-weather-visualizer-dev"
  location = var.gcp_region

  autogenerate_revision_name = false

  traffic {
    latest_revision = true
    percent         = 100
  }

  template {
    spec {
      container_concurrency = 5
      service_account_name  = google_service_account.kompal_weather_visualizer.email
      timeout_seconds       = 300

      containers {
        image = "asia.gcr.io/terraform-toshi0607/github_toshi0607_kompalweather/kompal-weather-visualizer-dev"
        env {
          name  = "GCS_BUCKET_NAME"
          value = google_storage_bucket.kompal_weather_report_dev.name
        }
        env {
          name  = "MAIL_SECRET_NAME"
          value = google_secret_manager_secret.google_user_email.secret_id
        }
        env {
          name  = "PW_SECRET_NAME"
          value = google_secret_manager_secret.google_password.secret_id
        }
        env {
          name  = "SERVICE_NAME"
          value = "kompal-weather-visualizer-dev"
        }
        env {
          name  = "VERSION"
          value = "0.1.0"
        }
        env {
          name  = "ENVIRONMENT"
          value = "development"
        }
        env {
          name  = "GCP_PROJECT_ID"
          value = var.gcp_project
        }
        env {
          name  = "SERVER_PORT"
          value = "8080"
        }

        resources {
          limits = {
            cpu    = "2000m"
            memory = "4096Mi"
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
    template.0.spec.0.containers.0.image]
  }

  depends_on = [
    google_project_iam_member.kompal_weather_visualizer_dev_is_logging_log_writer,
    google_project_iam_member.kompal_weather_visualizer_dev_is_secret_manager_admin
  ]
}
