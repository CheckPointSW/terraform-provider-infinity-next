terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.0.2"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_app_practice" "my-webapp-practice" {
  name = "some name"
  ips {
    performance_impact    = "VeryLow"    # enum of ["VeryLow", "LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "LowOrAbove" # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2020"       # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
    medium_confidence     = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
    low_confidence        = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
  }
  web_attacks {
    minimum_severity = "Critical" # enum of ["Critical", "High", "Medium"]
    advanced_setting {
      csrf_protection      = "Learn" # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      open_redirect        = "Learn" # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      error_disclosure     = "Learn" # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      body_size            = 1
      url_size             = 1
      header_size          = 1
      max_object_depth     = 1
      illegal_http_methods = false
    }
  }
  web_bot {
    inject_uris = ["url1", "url2"]
    valid_uris  = ["url1", "url2"]
  }
}
