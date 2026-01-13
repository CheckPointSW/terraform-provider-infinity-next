# Complete Web Application Asset Advanced Proxy Settings Configuration Example
#
# This example demonstrates ALL available configuration options for the
# inext_web_app_asset resource, including:
# - All 21 proxy settings
# - Additional NGINX instructions blocks (server and location)
# - Mutual TLS (mTLS) for both server and client
# - Custom headers
#

terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.4.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_app_asset" "compelte_proxy_settings_example" {
  name = "complete-web-app-example"
  urls = [
    "https://app.example.com",
    "https://www.example.com",
    "http://app.example.com", # HTTP version for redirect to HTTPS
    "http://www.example.com"  # HTTP version for redirect to HTTPS
  ]
  upstream_url = "https://backend.example.com"

  # -------------------------------------------------------------------------
  # ALL PROXY SETTINGS (Reverse Proxy Configuration)
  # -------------------------------------------------------------------------

  proxy_setting {
    key   = "connectTimeout"
    value = "60"
  }

  proxy_setting {
    key   = "readTimeout"
    value = "120"
  }

  proxy_setting {
    key   = "proxySendTimeout"
    value = "60"
  }

  proxy_setting {
    key   = "keepAliveTimeout"
    value = "75"
  }

  proxy_setting {
    key   = "proxySslName"
    value = "backend.example.com"
  }

  proxy_setting {
    key   = "proxySslVerify"
    value = "on"
  }

  proxy_setting {
    key   = "sslCiphersList"
    value = "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:!aNULL:!eNULL:!EXPORT:!DES:!MD5:!PSK:!RC4"
  }

  proxy_setting {
    key   = "activateHSTS"
    value = "true"
  }

  proxy_setting {
    key   = "enableHttp2"
    value = "true"
  }

  proxy_setting {
    key   = "enableWebSocketProxy"
    value = "true"
  }

  proxy_setting {
    key   = "proxyBufferSize"
    value = "8k"
  }

  proxy_setting {
    key   = "proxyBuffers"
    value = "16 8k"
  }

  proxy_setting {
    key   = "proxyBusyBuffersSize"
    value = "16k"
  }

  proxy_setting {
    key   = "healthCheck"
    value = "true"
  }

  proxy_setting {
    key   = "healthCheckPath"
    value = "/health"
  }

  proxy_setting {
    key   = "healthCheckInterval"
    value = "30"
  }

  proxy_setting {
    key   = "syslog"
    value = "true"
  }

  proxy_setting {
    key   = "enableProxyFaultLogging"
    value = "true"
  }

  proxy_setting {
    key   = "enableAppFaultLogging"
    value = "true"
  }

  proxy_setting {
    key   = "dnsServer"
    value = "8.8.8.8"
  }

  proxy_setting {
    key   = "proxyPassOverride"
    value = "https://internal-backend.example.com"
  }

  # -------------------------------------------------------------------------
  # Advanced Proxy Settings
  # -------------------------------------------------------------------------

  redirect_to_https = true
  access_log        = true
  custom_headers {
    name  = "X-Custom-Security-Token"
    value = "SecureToken123"
  }

  custom_headers {
    name  = "X-Application-Version"
    value = "v2.5.0"
  }

  custom_headers {
    name  = "X-Environment"
    value = "Production"
  }

  custom_headers {
    name  = "X-Request-Timestamp"
    value = "$time_iso8601"
  }

  # -------------------------------------------------------------------------
  # ADDITIONAL NGINX INSTRUCTIONS BLOCKS
  # -------------------------------------------------------------------------

  # Server-level NGINX configuration
  additional_instructions_blocks {
    type     = "server_instructions"
    filename = "server_config.conf"
    data     = file("${path.module}/nginx_server_instructions.conf")
    enable   = true
  }

  # Location-level NGINX configuration
  additional_instructions_blocks {
    type     = "location_instructions"
    filename = "location_config.conf"
    data     = file("${path.module}/nginx_location_instructions.conf")
    enable   = true
  }

  # -------------------------------------------------------------------------
  # MUTUAL TLS (mTLS) CONFIGURATION
  # -------------------------------------------------------------------------

  # Server-side mTLS certificate
  mtls {
    type             = "server"
    filename         = "server-cert.pem"
    certificate_type = ".pem"
    data             = file("${path.module}/certificates/server-certificate.pem")
    enable           = true
  }

  # Client-side mTLS certificate  
  mtls {
    type             = "client"
    filename         = "client-ca.pem"
    certificate_type = ".pem"
    data             = file("${path.module}/certificates/client-ca-bundle.pem")
    enable           = true
  }
}