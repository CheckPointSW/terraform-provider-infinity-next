terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.1.6"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_app_asset" "my-webapp-asset" {
  name         = "some name"
  profiles     = [inext_appsec_gateway_profile.my-appsec-gateway-profile.id, inext_docker_profile.my-docker-profile.id, inext_embedded_profile.my-embedded-profile.id, inext_kubernetes_profile.my-kubernetes-profile.id]
  behaviors    = [inext_trusted_sources.my-trusted-source-behavior.id, inext_exceptions.my-exceptions-behavior.id]
  upstream_url = "some url"
  urls         = ["http://host/path"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id       = inext_web_app_practice.my-webapp-practice.id # required
    triggers = [inext_log_trigger.mytrigger.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "HeaderKey" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value"]
  }
  tags {
    key   = "tagkey"
    value = "tagvalue"
  }
  mtls {
    filename         = "cert.der"
    certificate_type = ".der"
    data             = " cert data"
    type             = "client"
    enable           = true
  }
  additional_instructions_blocks {
    filename      = "location.json"
    filename_type = ".json"
    data          = "location data"
    type          = "location_instructions"
    enable        = true
  }
  redirect_to_https = "true"
  access_log        = "true"
  custom_headers {
    name  = "header1"
    value = "value1"
  }
}
