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

resource "inext_exceptions" "my-exceptions-behavior" {
  name = "some name"
  exception {
    match = { # currently matches with "AND" condition between all keys and their values
      host              = "www.acme.com"
      uri               = "/login"
      source_identifier = "value"
    }
    action  = "action1" # enum of ["drop", "skip", "accept", "suppressLog"]
    comment = "some comment"
  }
}
