resource "inext_log_trigger" "test" {
  name                             = "inext_log_trigger-test3"
  verbosity                        = "Extended" # enum of ["Minimal", "Standard", "Extended"]
  access_control_allow_events      = true
  access_control_drop_events       = true
  cef_ip_address                   = "10.0.0.1"
  cef_port                         = 81
  extend_logging                   = true
  extend_logging_min_severity      = "Critical" # enum of ["High", "Critical"]
  log_to_agent                     = true
  log_to_cef                       = true
  log_to_cloud                     = true
  log_to_syslog                    = true
  response_body                    = true
  response_code                    = true
  syslog_ip_address                = "10.0.0.2"
  syslog_port                      = 82
  threat_prevention_detect_events  = true
  threat_prevention_prevent_events = true
  web_body                         = true
  web_headers                      = false
  web_requests                     = true
  web_url_path                     = true
  web_url_query                    = true
}
