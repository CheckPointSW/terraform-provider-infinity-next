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

resource "inext_appsec_gateway_profile" "my-appsec-gateway-profile" {
  name                          = "my-appsec-gateway-profile"
  profile_sub_type              = "Azure"      # enum of ["Aws", "Azure", "VMware", "HyperV"]
  upgrade_mode                  = "Scheduled"  # enum of ["Automatic", "Manual", "Scheduled"]
  upgrade_time_schedule_type    = "DaysInWeek" # enum of ["DaysInMonth", "DaysInWeek", "Daily"]
  upgrade_time_hour             = "22:00"
  upgrade_time_duration         = 2
  upgrade_time_week_days        = ["Monday", "Thursday"]
  reverseproxy_upstream_timeout = 3600
  reverseproxy_additional_settings = {
    "Key1" = "Value"
    "Key2" = "Value2"
  }
  max_number_of_agents = 100
  additional_settings = {
    "Key1" = "Value"
    "Key2" = "Value2"
  }
  fail_open_inspection = true
  certificate_type     = "Vault" # enum of ["Vault", "Gateway"]
}
