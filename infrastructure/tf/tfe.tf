terraform {
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "credit-to-debit"

    workspaces {
      name = "credit-to-debit"
    }
  }
}