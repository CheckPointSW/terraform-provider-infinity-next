resource "inext_web_api_asset" "test" {
  name         = "inext_web_api_asset-test1"
  profiles     = [inext_appsec_gateway_profile.test2.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "some url 5"
  urls         = ["http://host5/path"]
  practice {
    main_mode = "Prevent" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      APIDiscovery = "Active"        # enum of ["Active", "Disabled"]
    }
    id       = inext_web_api_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }

  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  proxy_setting {
    key   = "another key"
    value = "another value"
  }
  proxy_setting {
    key   = "last key"
    value = "last value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie", "JWTKey"]
    values     = ["value3"]
  }
  source_identifier {
    identifier = "XForwardedFor" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie", "JWTKey"]
    values     = ["value2"]
  }
  source_identifier {
    identifier = "HeaderKey" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie", "JWTKey"]
    values     = ["value1"]
  }
  tags {
    key   = "tagkey1"
    value = "tagvalue1"
  }
  tags {
    key   = "tagkey2"
    value = "tagvalue2"
  }
  mtls {
    filename         = "cert.der"
    certificate_type = ".der"
    data             = "cert data"
    type             = "client"
    enable           = true
  }
}