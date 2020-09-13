// https://www.terraform.io/docs/providers/google/r/storage_bucket.html
resource "google_storage_bucket" "kompal_weather_report_dev" {
  name     = "kompal-weather-report-dev"
  location = "ASIA-NORTHEAST1"
}
