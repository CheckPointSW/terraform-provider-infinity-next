terraform {
  required_providers {
    inext = {
      version = "~> 1.0.0"
      source  = "checkpointsw/infinity-next"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_access_token" "access_token" {}
