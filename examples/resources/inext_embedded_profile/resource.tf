terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "~>1.3.2"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_embedded_profile" "my-embedded-profile" {
  name                       = "my-embedded-profile"
  max_number_of_agents       = 100
  upgrade_mode               = "Scheduled"  # enum of ["Automatic", "Manual", "Scheduled"]
  upgrade_time_schedule_type = "DaysInWeek" # enum of ["DaysInMonth", "DaysInWeek", "Daily"]
  upgrade_time_hour          = "22:00"
  upgrade_time_duration      = 2
  upgrade_time_week_days     = ["Monday", "Thursday"]
  only_defined_applications  = true
  additional_settings = {
    Key1 = "Value"
    Key2 = "Value2"
  }
}
