terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "~>1.2.0"
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
    match {
      operator = "and" # enum of ["and", "or", "not-equals", "equals", "in", "not-in", "exist"]
      operand {
        key   = "hostName" # enum of ["hostName", "sourceIdentifier", "url", "countryCode", "countryName", "manufacturer", "paramName", "paramValue", "protectionName", "sourceIP"]
        value = ["www.acme.com"]
      }
      operand {
        key   = "url"
        value = ["/login"]
      }
      operand {
        key   = "sourceIdentifier"
        value = ["value"]
      }
    }
    action  = "action1" # enum of ["drop", "skip", "accept", "suppressLog"]
    comment = "some comment"
  }
}
