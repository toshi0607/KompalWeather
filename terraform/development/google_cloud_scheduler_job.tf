// https://www.terraform.io/docs/providers/google/r/cloud_scheduler_job.html

locals {
  time_zone_tokyo = "Asia/Tokyo"
}

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
  schedule  = "*/15 0,14-23 * * mon,wed-sat"
  time_zone = local.time_zone_tokyo

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
  schedule  = "*/15 0,9-23 * * sun"
  time_zone = local.time_zone_tokyo

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
  schedule  = "*/15 0 * * tue"
  time_zone = local.time_zone_tokyo

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}

resource "google_cloud_scheduler_job" "daily_visualizer_dev" {
  attempt_deadline = "600s"
  name             = "daily-visualizer-dev"
  project          = var.gcp_project
  region           = var.gcp_region
  schedule         = "30 23 * * sun,mon,wed-sat"
  time_zone        = local.time_zone_tokyo

  retry_config {
    min_backoff_duration = "10s"
    max_doublings        = 2
    retry_count          = 3
  }

  http_target {
    body        = base64encode("{\"reportType\":\"daily\"}\n")
    http_method = "POST"
    uri         = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }
  }

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}

resource "google_cloud_scheduler_job" "weekly_visualizer_dev" {
  attempt_deadline = "600s"
  name             = "weekly-visualizer-dev"
  project          = var.gcp_project
  region           = var.gcp_region
  schedule         = "30 0 * * tue"
  time_zone        = local.time_zone_tokyo

  retry_config {
    min_backoff_duration = "10s"
    max_doublings        = 2
    retry_count          = 3
  }

  http_target {
    body        = base64encode("{\"reportType\":\"weekly\"}\n")
    http_method = "POST"
    uri         = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }
  }

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}

resource "google_cloud_scheduler_job" "monthly_visualizer_dev" {
  attempt_deadline = "600s"
  name             = "monthly-visualizer-dev"
  project          = var.gcp_project
  region           = var.gcp_region
  schedule         = "20 0 1 * *"
  time_zone        = local.time_zone_tokyo

  retry_config {
    min_backoff_duration = "10s"
    max_doublings        = 2
    retry_count          = 3
  }

  http_target {
    body        = base64encode("{\"reportType\":\"monthly\"}\n")
    http_method = "POST"
    uri         = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"

    oidc_token {
      audience              = "${google_cloud_run_service.kompal_weather_visualizer_dev.status.0.url}/visualize"
      service_account_email = google_service_account.kompal_weather_invoker.email
    }
  }

  depends_on = [google_project_iam_member.kompal_weather_invoker_dev_is_run_invoker]
}
