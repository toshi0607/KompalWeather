// https://www.terraform.io/docs/providers/google/r/secret_manager_secret.html
resource "google_secret_manager_secret" "kompal_url" {
  secret_id = "kompal_url"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "twitter_access_token" {
  secret_id = "twitter_access_token"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "twitter_access_token_secret" {
  secret_id = "twitter_access_token_secret"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "twitter_api_key" {
  secret_id = "twitter_api_key"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "twitter_api_key_secret" {
  secret_id = "twitter_api_key_secret"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "web_hook_url" {
  secret_id = "web_hook_url"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "google_password" {
  secret_id = "google_password"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "google_user_email" {
  secret_id = "google_user_email"
  replication {
    automatic = true
  }
}
