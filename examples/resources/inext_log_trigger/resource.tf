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

resource "inext_log_trigger" "mytrigger" {
  name                             = "mytrigger"
  acesss_control_allow_events      = false
  acesss_control_drop_events       = true
  threat_prevention_detect_events  = true
  threat_prevention_prevent_events = true
  web_body                         = false
  web_headers                      = false
  web_requests                     = false
  web_url_path                     = true
  web_url_query                    = true
  response_body                    = false
  response_code                    = true
  extend_logging                   = true
  extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
  log_to_agent                     = false
  log_to_cef                       = false
  cef_ip_address                   = "10.0.0.1"
  cef_port_num                     = 2000
  log_to_cloud                     = true
  log_to_syslog                    = true
  syslog_ip_address                = "10.10.10.10"
  syslog_port_num                  = 5004
  compliance_violations            = true
  compliance_warnings              = true
  verbosity                        = "Standard" # enum of ["Minimal", "Standard", "Extended"]
}
