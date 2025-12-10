terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "~>1.3.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_rate_limit_practice" "my_rate_limit_practice" {
  name       = "my rate limit practice"
  visibility = "Shared" # Optional: "Shared" (default) or "Local"

  # Multiple rate limit rules
  rule {
    uri     = "/api/v1/users"
    scope   = "Minute"              # Required: "Minute" or "Second"
    limit   = 100                   # Required: number of requests allowed
    action  = "Detect"              # Optional: "Detect", "Prevent", or "AccordingToPractice" (default)
    comment = "User API rate limit" # Optional: description of the rule
  }

  rule {
    uri     = "/api/v1/login"
    scope   = "Second"
    limit   = 5
    action  = "Prevent"
    comment = "Login endpoint strict rate limit"
  }

  rule {
    uri     = "/api/v1/search"
    scope   = "Minute"
    limit   = 200
    action  = "AccordingToPractice"
    comment = "Search API with higher limit"
  }
}

# Example with minimal configuration (only required fields)
resource "inext_rate_limit_practice" "minimal_example" {
  name = "minimal rate limit practice"

  rule {
    uri   = "/api/minimal"
    scope = "Minute"
    limit = 50
  }
}