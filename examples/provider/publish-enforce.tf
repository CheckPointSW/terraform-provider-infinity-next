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

variable "profile_ids" {
  type        = list(string)
  default     = []
  description = "List of profile IDs to enforce. If empty, all profiles will be enforced."
}

# Publish and Enforce resource - should be the last resource to be applied
# Only ONE instance is allowed per provider/account
resource "inext_publish_enforce" "publish-and-enforce" {
  publish     = var.publish
  enforce     = var.enforce
  profile_ids = var.profile_ids

  # This resource should depend on all other resources
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
