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

resource "inext_web_api_practice" "my-webapi-practice" {
  name = "some name"
  ips {
    performance_impact    = "MediumOrLower" # enum of ["LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "LowOrAbove"    # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2020"          # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "Prevent"       # enum of ["Detect", "Prevent", "Inactive"]
    medium_confidence     = "Detect"        # enum of ["Detect", "Prevent", "Inactive"]
    low_confidence        = "Inactive"      # enum of ["Detect", "Prevent", "Inactive"]
  }
  api_attacks {
    minimum_severity = "Critical" # enum of ["Critical", "High", "Medium"]
    advanced_setting {
      body_size            = 1000
      url_size             = 1000
      header_size          = 1000
      max_object_depth     = 1000
      illegal_http_methods = true
    }
  }
  schema_validation {
    filename = basename(data.local_file.schema_validation_file.filename)
    data     = data.local_file.schema_validation_file.content
  }
}
