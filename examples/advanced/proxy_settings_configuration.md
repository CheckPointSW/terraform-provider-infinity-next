# Proxy Settings Configuration Guide

## Overview

Proxy settings allow you to configure advanced reverse proxy behavior for your web application asset. These settings are passed as key-value pairs in the `proxy_setting` block of the `web-api/app-asset` resource.

CloudGuard WAF (Infinity Next) uses NGINX as its underlying reverse proxy engine. These proxy settings correspond to NGINX configuration directives that control how the AppSec Gateway handles traffic between clients and your protected upstream servers.

## Related Documentation

- [CloudGuard WAF - Edit Reverse Proxy Advanced Settings](https://waf-doc.inext.checkpoint.com/how-to/edit-reverse-proxy-advanced-settings-for-a-web-asset)
- [NGINX Documentation](https://nginx.org/en/docs/)
- [NGINX HTTP Proxy Module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)

## Configuration

Proxy settings are configured using the `proxy_setting` block within your web api/app asset resource:

```hcl
resource "infinitynext_web_app_asset" "example" {
  name = "my-app"
  urls = ["https://example.com"]
  
  proxy_setting {
    key   = "connectTimeout"
    value = "60"
  }
  
  proxy_setting {
    key   = "enableHttp2"
    value = "true"
  }
  
  # Add more proxy settings as needed
}
```

## Available Proxy Settings

### Timeout Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `connectTimeout` | Connect Timeout | `"60"` | Defines a timeout to establish a connection with the proxied server. Maps to NGINX `proxy_connect_timeout` directive. |
| `keepAliveTimeout` | Keep Alive Timeout | `"75"` | Sets a timeout in which a keep-alive client connection stays open on the server side. The zero value disables keep-alive client connections. Maps to NGINX `keepalive_timeout` directive. |
| `proxySendTimeout` | Proxy Send Timeout | `"60"` | Sets a timeout to transmit a request to the proxied server. The timeout is set only between two successive write operations, not for the transmission of the whole request. If the proxied server does not receive anything within this time, the connection is closed. Maps to NGINX `proxy_send_timeout` directive. |
| `readTimeout` | Read Timeout | `"60"` | Defines a timeout to read a response from the proxied server. The timeout is set only between two successive read operations, not for the transmission of the whole response. If the proxied server does not transmit anything within this time, the connection is closed. Maps to NGINX `proxy_read_timeout` directive. |

### SSL/TLS Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `proxySslName` | Proxy SSL Name | `"example.com"` | Allows overriding the server name used to verify the certificate of the proxied HTTPS server, and allows it to pass through SNI when establishing a connection with the proxied HTTPS server. Maps to NGINX `proxy_ssl_name` directive. |
| `proxySslVerify` | Proxy SSL Verify | `"on"` or `"off"` | Enables or disables verification of the proxied HTTPS server certificate. Best practice is to enable this in production environments. Maps to NGINX `proxy_ssl_verify` directive. |
| `sslCiphersList` | SSL Ciphers - separated by ':' | `"HIGH:!aNULL:!MD5"` | Allows the definition of an additional SSL ciphers' list that will be allowed by the CloudGuard WAF AppSec Gateway for HTTPS traffic. Ciphers are delimited by ':'. Maps to NGINX `ssl_ciphers` directive. See [OpenSSL Cipher List Format](https://www.openssl.org/docs/man1.1.1/man1/ciphers.html) and [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/) for recommended cipher suites. |
| `activateHSTS` | Enable HSTS | `"true"` or `"false"` | Controls the activation of Strict Transport Security Response header on all responses coming from the CloudGuard WAF AppSec Gateway. False by default. When enabled, adds the Strict-Transport-Security header to responses. |

### Protocol and Feature Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `enableHttp2` | Enable HTTP/2 | `"true"` or `"false"` | Controls the activation of HTTP/2 protocol support. See [NGINX HTTP/2 Module](https://nginx.org/en/docs/http/ngx_http_v2_module.html) for details. |
| `enableWebSocketProxy` | Enable WebSocket proxying | `"true"` or `"false"` | Controls the activation of a proxy tunnel for WebSocket requests through the CloudGuard WAF AppSec gateway. True by default. See [NGINX WebSocket proxying](https://nginx.org/en/docs/http/websocket.html) for implementation details. |

### Buffer Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `proxyBufferSize` | Proxy Buffer Size | `"4k"` or `"8k"` | Sets the size of the buffer used for reading the first part of the response received from the proxied server (usually contains the response header). Maps to NGINX `proxy_buffer_size` directive. |
| `proxyBuffers` | Proxy Buffers | `"8 4k"` | Sets the number and size of the buffers used for reading a response from the proxied server. Format: `number size`. For example, "8 4k" means 8 buffers of 4 kilobytes each. Maps to NGINX `proxy_buffers` directive. |
| `proxyBusyBuffersSize` | Proxy Busy Buffers Size | `"8k"` | When buffering of responses from the proxied server is enabled, limits the total size of buffers that can be busy sending a response to the client while the response is not yet fully read. Maps to NGINX `proxy_busy_buffers_size` directive. |

### Health Check Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `healthCheck` | Health Check | `"true"` or `"false"` | Controls periodic health checks to verify NGINX is running. True by default. |
| `healthCheckPath` | Health Check Path for SaaS | `"/health"` | **Note:** This setting is effective for CloudGuard WAF as a Service only. Set the URL path for the health check request. This path verifies the service's availability and responsiveness. Must start with "/". |
| `healthCheckInterval` | Health Check Interval for SaaS (in Sec) | `"30"` | **Note:** This setting is effective for CloudGuard WAF as a Service only. Indicates how often (in seconds) the health check request is sent to the specified path. A lower number means more frequent monitoring, while a higher number reduces the check frequency. Valid range: 30 to 3600 seconds. |

### Logging Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `syslog` | Enable Failed Requests Logging | `"true"` or `"false"` | When set to "true", the CloudGuard WAF AppSec Gateway will also log all blocked traffic due to Reverse Proxy configuration. False by default. Logs are sent to syslog. |
| `enableProxyFaultLogging` | Reverse Proxy Error Logging | `"true"` or `"false"` | Enable logging of reverse proxy errors. Helps identify issues at the reverse proxy level such as connection failures, timeout errors, or SSL/TLS issues. |
| `enableAppFaultLogging` | Upstream Application Error Logging | `"true"` or `"false"` | Enable logging of upstream application errors. Helps identify issues with the protected backend server such as 5xx errors or application-level failures. |

### Other Settings

| Key | Description | Example Value | Notes |
|-----|-------------|---------------|-------|
| `dnsServer` | DNS Server | `"8.8.8.8"` or `"1.1.1.1"` | Configures a domain name server's IP address used to resolve names of upstream URLs (the backend servers' URL) into addresses. Useful when your upstream servers use domain names that need custom DNS resolution. Maps to NGINX `resolver` directive. |
| `proxyPassOverride` | Upstream URL Override | `"https://internal.example.com"` | Allows overriding the upstream URL for specific routing needs. This can be useful for directing traffic to a different backend server than the one specified in the `upstream_url` attribute. Maps to NGINX `proxy_pass` directive. |

## Complete Example

```hcl
resource "infinitynext_web_app_asset" "example" {
  name         = "production-app"
  urls         = ["https://app.example.com"]
  upstream_url = "https://backend.example.com"
  state        = "Active"

  # Timeout configurations
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

  # Enable HTTP/2 and WebSockets
  proxy_setting {
    key   = "enableHttp2"
    value = "true"
  }
  
  proxy_setting {
    key   = "enableWebSocketProxy"
    value = "true"
  }

  # SSL/TLS configuration
  proxy_setting {
    key   = "activateHSTS"
    value = "true"
  }
  
  proxy_setting {
    key   = "proxySslVerify"
    value = "on"
  }
  
  proxy_setting {
    key   = "proxySslName"
    value = "backend.example.com"
  }
  
  proxy_setting {
    key   = "sslCiphersList"
    value = "ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:HIGH:!aNULL:!MD5"
  }

  # Buffer settings for large responses
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

  # Health check configuration
  proxy_setting {
    key   = "healthCheck"
    value = "true"
  }
  
  proxy_setting {
    key   = "healthCheckPath"
    value = "/api/health"
  }
  
  proxy_setting {
    key   = "healthCheckInterval"
    value = "30"
  }

  # Logging configuration
  proxy_setting {
    key   = "enableProxyFaultLogging"
    value = "true"
  }
  
  proxy_setting {
    key   = "enableAppFaultLogging"
    value = "true"
  }
  
  proxy_setting {
    key   = "syslog"
    value = "true"
  }

  # Custom DNS
  proxy_setting {
    key   = "dnsServer"
    value = "8.8.8.8"
  }

  # Upstream URL override for specific routing
  proxy_setting {
    key   = "proxyPassOverride"
    value = "https://internal-backend.example.com"
  }
}
```

### NGINX Configuration Mapping

The proxy settings you configure via Terraform map directly to NGINX directives. For detailed information about each NGINX directive, refer to:

- [NGINX HTTP Proxy Module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html) - For proxy-related settings
- [NGINX HTTP SSL Module](https://nginx.org/en/docs/http/ngx_http_ssl_module.html) - For SSL/TLS settings
- [NGINX HTTP Core Module](https://nginx.org/en/docs/http/ngx_http_core_module.html) - For core HTTP settings

### Additional NGINX Configuration

For more advanced NGINX configuration beyond the available proxy settings, you can use the `additional_instructions_blocks` attribute in the web api/app asset resource. This allows you to upload custom NGINX configuration snippets that will be inserted into:
- **Location blocks** (`location_instructions`) - For request-specific settings
- **Server blocks** (`server_instructions`) - For server-level settings

See the [complete_proxy_settings.tf example](complete_proxy_settings.tf) for details on using `additional_instructions_blocks`.

**Example files and comprehensive usage guide:**
- [Server Instructions Example](nginx_server_instructions.conf) - Server-level NGINX directives
- [Location Instructions Example](nginx_location_instructions.conf) - Location-level NGINX directives

## Additional Resources

### CloudGuard WAF Documentation
- [CloudGuard WAF Documentation](https://waf-doc.inext.checkpoint.com/)
- [Edit Reverse Proxy Advanced Settings for a Web Asset](https://waf-doc.inext.checkpoint.com/how-to/edit-reverse-proxy-advanced-settings-for-a-web-asset)

### NGINX Documentation
- [NGINX Official Documentation](https://nginx.org/en/docs/)
- [NGINX HTTP Proxy Module Reference](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [NGINX HTTP SSL Module](https://nginx.org/en/docs/http/ngx_http_ssl_module.html)
- [WebSocket Proxying with NGINX](https://nginx.org/en/docs/http/websocket.html)
- [Configuring HTTPS Servers](https://nginx.org/en/docs/http/configuring_https_servers.html)

### SSL/TLS Cipher Configuration
- [OpenSSL Cipher List Format](https://www.openssl.org/docs/man1.1.1/man1/ciphers.html) - Official OpenSSL cipher string format and syntax
