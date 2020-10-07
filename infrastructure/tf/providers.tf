provider "google" {
  credentials = var.gcp_key 
  project     = "my-project-id"
  region      = "us-central1"
}