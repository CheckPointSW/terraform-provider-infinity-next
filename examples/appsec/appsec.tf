terraform {
  required_providers {
    inext = {
      source = "CheckPointSW/infinity-next"
      version = "1.1.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_app_asset" "prod_eu_acme" {
  name         = "ACME Power"
  profiles     = [inext_appsec_gateway_profile.eu_appsec.id]
  behaviors    = [inext_trusted_sources.new_parameter.id, inext_exceptions.shared_region_exceptions.id]
  upstream_url = "https://www.acme.com"
  urls         = ["http://192.168.2.200/main"]
  practice {
    main_mode = "Prevent" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS              = "AccordingToPractice"
      APIAttacks       = "AccordingToPractice"
      SchemaValidation = "Disabled"
      Snort            = "AccordingToPractice"
    }
    id       = inext_web_app_practice.eu_acme_protection.id # required
    triggers = [inext_log_trigger.log_trigger.id]
  }
  source_identifier {
    identifier = "HeaderKey"
    values     = ["login"]
  }
}

resource "inext_web_app_practice" "eu_acme_protection" {
  name = "EU ACME API Protection"
  ips {
    performance_impact    = "MediumOrLower" # enum of ["VeryLow", "LowOrLower", "MediumOrLower", "HighOrLower"]
    severity_level        = "MediumOrAbove" # enum of ["LowOrAbove", "MediumOrAbove", "HighOrAbove", "Critical"]
    protections_from_year = "2016"          # enum of ["1999", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"]
    high_confidence       = "Prevent"       # enum of ["Detect", "Prevent", "Inactive"]
    medium_confidence     = "Prevent"       # enum of ["Detect", "Prevent", "Inactive"]
    low_confidence        = "Detect"        # enum of ["Detect", "Prevent", "Inactive"]
  }
  web_attacks {
    minimum_severity = "High" # enum of ["Critical", "High", "Medium"]
    advanced_setting {
      csrf_protection      = "Prevent"             # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      open_redirect        = "Disabled"            # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      error_disclosure     = "AccordingToPractice" # enum of ["Disabled", "Learn", "Prevent", "AccordingToPractice"]
      body_size            = 1000000
      url_size             = 32768
      header_size          = 102400
      max_object_depth     = 100
      illegal_http_methods = false
    }
  }
}

resource "inext_exceptions" "shared_region_exceptions" {
  name = "Shared Region Exceptions"
  exception {
    match = {
      url              = "/"
      sourceIdentifier = "1.1.1.1/24"
    }
    action  = "drop"
    comment = "some comment"
  }
  exception {
    match = {
      url           = "/"
      parameterName = "token"
    }
    action = "skip"
  }
}

resource "inext_trusted_sources" "new_parameter" {
  name               = "Trusted Sources"
  min_num_of_sources = 3
  sources_identifiers = [
    "ruth@acme.com",
    "mark@acme.com",
    "cathy@acme.com",
    "alice@acme.com",
    "joel@acme.com",
    "ben@acme.com"
  ]
}

resource "inext_log_trigger" "log_trigger" {
  verbosity                        = "Standard"
  access_control_allow_events      = false
  access_control_drop_events       = true
  extend_logging                   = true
  extend_logging_min_severity      = "High"
  log_to_agent                     = false
  log_to_cef                       = false
  log_to_cloud                     = true
  log_to_syslog                    = false
  name                             = "AppSec Log Trigger"
  response_body                    = true
  response_code                    = false
  threat_prevention_detect_events  = true
  threat_prevention_prevent_events = true
  web_body                         = false
  web_headers                      = false
  web_requests                     = false
  web_url_path                     = true
  web_url_query                    = true
}

resource "inext_appsec_gateway_profile" "eu_appsec" {
  name                 = "EU AppSec Gateways"
  profile_sub_type     = "Azure"
  upgrade_mode         = "Automatic"
  max_number_of_agents = 100
}
