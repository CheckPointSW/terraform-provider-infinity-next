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

resource "inext_web_api_asset" "my-webapi-asset" {
  name            = "some name"
  profiles        = [inext_appsec_gateway_profile.my-appsec-gateway-profile.id, inext_docker_profile.my-docker-profile.id, inext_embedded_profile.my-embedded-profile.id, inext_kubernetes_profile.my-kubernetes-profile.id]
  trusted_sources = [inext_trusted_sources.my-trusted-source-behavior.id]
  upstream_url    = "some url"
  urls            = ["some url"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_api_practice.my-webapi-practice.id # required
    triggers   = [inext_log_trigger.mytrigger.id]
    exceptions = [inext_exceptions.my-exceptions-behavior.id]
  }
  proxy_setting {
    key   = "some key"
    value = "some value"
  }
  source_identifier {
    identifier = "XForwardedFor" # enum of ["SourceIP", "XForwardedFor", "HeaderKey", "Cookie"]
    values     = ["value1", "value2"]
  }
}
