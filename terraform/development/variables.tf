variable "gcp_project" {
  type        = string
  description = "GCP project name"
}

# https://cloud.google.com/compute/docs/regions-zones
variable "gcp_region" {
  type        = string
  description = "GCP region"
}

variable "gcp_zone" {
  type        = string
  description = "GCP zone"
}

variable "tfstate_bucket" {
  type = string
  description = "GCS bucket name to store tfstate"
}

variable "scheduler_target" {
  type = string
  description = "Audience endpoint of scheduler"
}
