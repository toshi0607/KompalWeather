// https://www.terraform.io/docs/providers/google/r/cloud_scheduler_job.html
resource "google_cloud_scheduler_job" "dev1" {
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }

    uri = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
  }

  name      = "kompal-weather-dev"
  project   = var.gcp_project
  region    = var.gcp_region
  schedule  = " */15 0,15-23 * * mon,wed-sat"
  time_zone = "Asia/Tokyo"

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}

resource "google_cloud_scheduler_job" "dev2" {
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }

    uri = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
  }

  name      = "kompal-weather-dev2"
  project   = var.gcp_project
  region    = var.gcp_region
  schedule  = "*/15 0,10-23 * * sun"
  time_zone = "Asia/Tokyo"

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}

resource "google_cloud_scheduler_job" "dev3" {
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }

    uri = "${google_cloud_run_service.kompal-weather-dev.status.0.url}/watch"
  }

  name      = "kompal-weather-dev3"
  project   = var.gcp_project
  region    = var.gcp_region
  schedule  = "00 0 * * tue"
  time_zone = "Asia/Tokyo"

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}
