terraform {
  required_providers {
    inext = {
      version = "~> 1.0.0"
      source  = "checkpointsw/infinitynext"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_docker_profile" "my-docker-profile" {
  name                      = "my-docker-profile"
  max_number_of_agents      = 100
  only_defined_applications = true
  additional_settings = {
    Key1 = "Value"
    Key2 = "Value2"
  }
}
