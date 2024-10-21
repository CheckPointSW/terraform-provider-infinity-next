resource "inext_web_app_asset" "test1" {
  name         = "inext_web_app_asset-test1"
  profiles     = [inext_appsec_gateway_profile.test.id]
  behaviors    = [inext_trusted_sources.test.id, inext_exceptions.test.id]
  upstream_url = "some url5"
  urls         = ["http://host/path5"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
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

resource "inext_web_app_asset" "test2" {
  name            = "inext_web_app_asset-test2"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url6"
  urls            = ["http://host/path6"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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

resource "inext_web_app_asset" "test3" {
  name            = "inext_web_app_asset-test3"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url7"
  urls            = ["http://host/path7"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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
  name            = "inext_web_app_asset-test4"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url8"
  urls            = ["http://host/path8"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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
  name            = "inext_web_app_asset-test5"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url9"
  urls            = ["http://host/path9"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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
  name            = "inext_web_app_asset-test6"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url10"
  urls            = ["http://host/path10"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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
  name            = "inext_web_app_asset-test7"
  profiles        = [inext_appsec_gateway_profile.test.id]
  trusted_sources = [inext_trusted_sources.test.id]
  upstream_url    = "some url11"
  urls            = ["http://host/path11"]
  practice {
    main_mode = "Learn" # enum of ["Prevent", "Inactive", "Disabled", "Learn"]
    sub_practices_modes = {
      IPS    = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      WebBot = "AccordingToPractice" # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
      Snort  = "Disabled"            # enum of ["Detect", "Prevent", "Inactive", "AccordingToPractice", "Disabled", "Learn", "Active"]
    }
    id         = inext_web_app_practice.test.id # required
    triggers   = [inext_log_trigger.test.id]
    exceptions = [inext_exceptions.test.id]
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