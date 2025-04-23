resource "inext_web_app_practice" "test" {
  name = "inext_web_app_practice-test1"
  ips {
    performance_impact    = "LowOrLower" # enum of ["LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "LowOrAbove" # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2016"       # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "Detect"     # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    medium_confidence     = "Detect"     # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    low_confidence        = "Detect"     # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
  }
  web_attacks {
    minimum_severity = "High" # enum of ["Critical", "High", "Medium"]
    advanced_setting {
      csrf_protection      = "Prevent"             # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      open_redirect        = "Disabled"            # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      error_disclosure     = "AccordingToPractice" # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      body_size            = 1000
      url_size             = 1000
      header_size          = 1000
      max_object_depth     = 1000
      illegal_http_methods = true
    }
  }
  web_bot {
    inject_uris = ["url1", "url2"]
    valid_uris  = ["url1", "url2"]
  }
  file_security {
    severity_level               = "LowOrAbove"          # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    high_confidence              = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    medium_confidence            = "Detect"              # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    low_confidence               = "Inactive"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_file_size_limit        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    file_size_limit              = 10
    file_size_limit_unit         = "MB"                  # enum of ["Bytes","KB", "MB", "GB"]
    files_without_name           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    required_archive_extraction  = true
    archive_file_size_limit      = 100
    archive_file_size_limit_unit = "MB"                  # enum of ["Bytes","KB", "MB", "GB"]
    allow_archive_within_archive = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_an_unopened_archive    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice"]
    allow_file_type              = true
    required_threat_emulation    = true
  }
}