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
          value = "terraform-toshi0607"
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
          value = "web_hook_url"
        }
        env {
          name  = "KOMPAL_URL_SECRET_NAME"
          value = "kompal_url"
        }
        env {
          name  = "TWITTER_ACCESS_TOKEN_SECRET_NAME"
          value = "twitter_access_token"
        }
        env {
          name  = "TWITTER_ACCESS_TOKEN_SECRET_SECRET_NAME"
          value = "twitter_access_token_secret"
        }
        env {
          name  = "TWITTER_API_KEY_SECRET_NAME"
          value = "twitter_api_key"
        }
        env {
          name  = "TWITTER_API_KEY_SECRET_SECRET_NAME"
          value = "twitter_api_key_secret"
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
    google_project_iam_member.kompal_weather_is_logging_logwriter,
    google_project_iam_member.kompal_weather_is_secretmanager_admin
  ]
}
