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
    inext_web_app_practice.my-webapp-practice,
    inext_web_api_asset.my-webapi-asset,
    inext_web_api_practice.my-webapi-practice,
    inext_appsec_gateway_profile.my-appsec-gateway-profile,
    inext_log_trigger.mytrigger,
    inext_exceptions.my-exceptions-behavior,
    inext_trusted_sources.my-trusted-source-behavior,
    inext_rate_limit_practice.my-rate-limit-practice,
  ]
}
