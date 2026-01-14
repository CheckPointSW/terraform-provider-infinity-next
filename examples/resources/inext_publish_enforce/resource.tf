terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "~>1.5.0"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

variable "publish" {
  type        = bool
  default     = false
  description = "Set to true to trigger publish operation"
}

variable "enforce" {
  type        = bool
  default     = false
  description = "Set to true to trigger enforce operation"
}

# Example resources that would be created before publish/enforce
resource "inext_web_app_asset" "my-webapp-asset" {
  name = "my web app"
  urls = ["http://example.com"]
  # ... other configuration
}

resource "inext_appsec_gateway_profile" "my-appsec-gateway-profile" {
  name = "my appsec gateway profile"
  # ... other configuration
}

# Publish and Enforce resource - should be the last resource to be applied
# Only ONE instance is allowed per provider/account
resource "inext_publish_enforce" "publish-and-enforce" {
  publish = var.publish
  enforce = var.enforce

  # Optional: specify profile IDs to enforce directly in the resource
  # If empty or not provided, all profiles will be enforced
  # profile_ids = ["profile-id-1", "profile-id-2"]

  # IMPORTANT: depends_on MUST include ALL other resources to ensure
  # publish/enforce runs last and avoids conflicts
  depends_on = [
    inext_web_app_asset.my-webapp-asset,
    inext_appsec_gateway_profile.my-appsec-gateway-profile,
    # Add all other resources that should be created before publish/enforce
  ]
}
