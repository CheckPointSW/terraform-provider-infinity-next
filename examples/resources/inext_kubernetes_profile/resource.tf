terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "~>1.4.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_kubernetes_profile" "my-kubernetes-profile" {
  name                      = "my-kubernetes-profile"
  profile_sub_type          = "AppSec"
  max_number_of_agents      = 100
  only_defined_applications = true
  additional_settings = {
    Key1 = "Value"
    Key2 = "Value2"
  }
}
