terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.5.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_trusted_sources" "my-trusted-source-behavior" {
  name                = "some name"
  visibility          = "Shared"
  min_num_of_sources  = 1
  sources_identifiers = ["identifier1", "identifier2"]
}
