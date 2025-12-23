resource "inext_web_app_asset" "test1" {
  name         = "inext_web_app_asset-test1"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url5.com"
  urls         = ["http://host/path5"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
  tags {
    key   = "tagkey"
    value = "tagvalue"
  }
  mtls {
    filename         = "cert.der"
    certificate_type = ".der"
    data             = file("${path.module}/cert.der") # file content - change path to your file
    type             = "client"                        # enum of ["client", "server"]
    enable           = true
  }
  additional_instructions_blocks {
    filename      = "location.json"
    filename_type = ".json"
    data          = file("${path.module}/location.json") # file content - change path to your file
    type          = "location_instructions"              # enum of ["location_instructions", "server_instructions"]
    enable        = true
  }
  redirect_to_https = "true"
  access_log        = "true"
  custom_headers {
    name  = "header1"
    value = "value1"
  }
  is_shares_urls = "true"
}

resource "inext_web_app_asset" "test-server-instructions" {
  name         = "inext_web_app_asset-test-server-instructions"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url5.com"
  urls         = ["http://host/path5"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
  tags {
    key   = "tagkey"
    value = "tagvalue"
  }
  mtls {
    filename         = "cert.der"
    certificate_type = ".der"
    data             = file("${path.module}/cert.der") # file content - change path to your file
    type             = "client"                        # enum of ["client", "server"]
    enable           = true
  }
  additional_instructions_blocks {
    filename      = "server.json"
    filename_type = ".json"
    data          = file("${path.module}/server.json") # file content - change path to your file
    type          = "server_instructions"              # enum of ["location_instructions", "server_instructions"]
    enable        = true
  }
  redirect_to_https = "true"
  access_log        = "true"
  custom_headers {
    name  = "header1"
    value = "value1"
  }
  is_shares_urls = "true"
}

resource "inext_web_app_asset" "test2" {
  name         = "inext_web_app_asset-test2"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "https://some url6.com"
  urls         = ["http://host/path6"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
  tags {
    key   = "tagkey"
    value = "tagvalue"
  }
  mtls {
    filename         = "cert2.der"
    certificate_type = ".der"
    data             = file("${path.module}/cert.der") # file content - change path to your file
    type             = "server"                        # enum of ["client", "server"]
    enable           = true
  }
}

resource "inext_web_app_asset" "test3" {
  name         = "inext_web_app_asset-test3"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url7.com"
  urls         = ["http://host/path7"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
}

resource "inext_web_app_asset" "test4" {
  name         = "inext_web_app_asset-test4"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url8.com"
  urls         = ["http://host/path8"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
}

resource "inext_web_app_asset" "test5" {
  name         = "inext_web_app_asset-test5"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "some url9"
  urls         = ["http://host/path9"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
}

resource "inext_web_app_asset" "test6" {
  name         = "inext_web_app_asset-test6"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url10.com"
  urls         = ["http://host/path10"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
}

resource "inext_web_app_asset" "test7" {
  name         = "inext_web_app_asset-test7"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "http://some url11.com"
  urls         = ["http://host/path11"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS           = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot        = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      FileSecurity  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.test.id # required
    triggers = [inext_log_trigger.test.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "SourceIP" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
}