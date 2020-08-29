// https://www.terraform.io/docs/providers/google/r/cloud_scheduler_job.html
resource "google_cloud_scheduler_job" "dev1" {
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"

    oidc_token {
      audience              = var.scheduler_target
      service_account_email = google_service_account.kompal_weather_invoker.email
    }

    // It's better to refer to cloud run resource output
    uri = var.scheduler_target
  }

  name      = "kompal-weather-dev"
  project   = var.gcp_project
  region    = var.gcp_region
  schedule  = " */15 0,15-23 * * mon,wed-sat"
  time_zone = "Asia/Tokyo"
}

resource "google_cloud_scheduler_job" "dev2" {
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"

    oidc_token {
      audience              = var.scheduler_target
      service_account_email = google_service_account.kompal_weather_invoker.email
    }

    uri = var.scheduler_target
  }

  name      = "kompal-weather-dev2"
  project   = var.gcp_project
  region    = var.gcp_region
  schedule  = "*/15 0,10-23 * * sun"
  time_zone = "Asia/Tokyo"
}
