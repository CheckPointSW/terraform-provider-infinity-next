resource "inext_web_app_practice" "test" {
  name = "inext_web_app_practice-test1"
  ips {
    performance_impact    = "LowOrLower" # enum of ["LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "LowOrAbove" # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2016"       # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
    medium_confidence     = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
    low_confidence        = "Detect"     # enum of ["Detect", "Prevent", "Inactive"]
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
}