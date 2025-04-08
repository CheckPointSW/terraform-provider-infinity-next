terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.1.1"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_api_practice" "my-webapi-practice" {
  name       = "some name"
  visibility = "Shared" # enum of ["Shared", "Local"]
  ips {
    performance_impact    = "MediumOrLower"       # enum of ["LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "LowOrAbove"          # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2020"                # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    medium_confidence     = "Detect"              # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    low_confidence        = "Inactive"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
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
  file_security {
    severity_level               = "LowOrAbove"          # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    high_confidence              = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    medium_confidence            = "Detect"              # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    low_confidence               = "Inactive"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_file_size_limit        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    file_size_limit              = 10
    file_size_limit_unit         = "MB"                  # enum of ["Bytes","KB", "MB", "GB"]
    files_without_name            = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    required_archive_extraction  = true
    archive_file_size_limit      = 100
    archive_file_size_limit_unit = "MB"                  # enum of ["Bytes","KB", "MB", "GB"]
    allow_archive_within_archive = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_an_unopened_archive    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_file_type              = true
    required_threat_emulation    = true
  }
}
