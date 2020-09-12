// https://www.terraform.io/docs/providers/google/r/google_project_iam.html
resource "google_project_iam_member" "kompal_weather_is_logging_logwriter" {
  member  = "serviceAccount:${google_service_account.kompal_weather.email}"
  project = var.gcp_project
  role    = "roles/logging.logWriter"
}

resource "google_project_iam_member" "kompal_weather_is_secretmanager_admin" {
  member  = "serviceAccount:${google_service_account.kompal_weather.email}"
  project = var.gcp_project
  role    = "roles/secretmanager.admin"
}

resource "google_project_iam_member" "kompal_weather_is_monitoring_editor" {
  member  = "serviceAccount:${google_service_account.kompal_weather.email}"
  project = var.gcp_project
  role    = "roles/monitoring.editor"
}

resource "google_project_iam_member" "kompal_weather_invoker_dev_is_run_invoker" {
  member  = "serviceAccount:${google_service_account.kompal_weather_invoker.email}"
  project = var.gcp_project
  role    = "roles/run.invoker"
}

resource "google_project_iam_member" "kompal_weather_visualizer_dev_is_storage_object_admin" {
  member  = "serviceAccount:${google_service_account.kompal_weather_visualizer.email}"
  project = var.gcp_project
  role    = "roles/storage.objectAdmin"
}

resource "google_project_iam_member" "kompal_weather_visualizer_dev_is_secretmanager_admin" {
  member  = "serviceAccount:${google_service_account.kompal_weather_visualizer.email}"
  project = var.gcp_project
  role    = "roles/secretmanager.admin"
}

resource "google_project_iam_member" "kompal_weather_visualizer_dev_is_logging_logwriter" {
  member  = "serviceAccount:${google_service_account.kompal_weather_visualizer.email}"
  project = var.gcp_project
  role    = "roles/logging.logWriter"
}
