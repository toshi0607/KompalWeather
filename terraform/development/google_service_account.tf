// https://www.terraform.io/docs/providers/google/r/google_service_account.html
resource "google_service_account" "kompal_weather" {
  account_id   = "kompal-weather"
  display_name = "kompal-weather"
  project      = var.gcp_project
}

resource "google_service_account" "kompal_weather_invoker" {
  account_id   = "kompal-weather-invoker-dev"
  display_name = "kompal-weather-invoker-dev"
  project      = var.gcp_project
}
