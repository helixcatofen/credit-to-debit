resource "google_storage_bucket" "bucket" {
  name = "credit-to-debit-bucket"
}


resource "google_cloudfunctions_function" "function" {
  name        = "function-test"
  description = "My function"
  runtime     = "go113"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  entry_point           = "helloGET"
  environment_variables = {
      STARLING_API_KEY = var.starling_api_key
      TRUELAYER_ID = var.truelayer_id
      TRUELAYER_SECRET = var.truelayer_secret
      TRUELAYER_TOKEN = var.truelayer_token
  }
}

# IAM entry for all users to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}