provider "google" {
  credentials = var.gcp_key 
  project     = var.name
  region      = var.region
}