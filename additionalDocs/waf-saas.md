# WAF SaaS — Terraform Provider Status

This document describes the current status of WAF SaaS support in this Terraform provider, including the deployment flow, known limitations, and technical gaps.

## Deployment Flow

### 1. Create the SaaS Profile in the UI

WAF SaaS profile creation is not yet supported via Terraform. You must create the AppSec SaaS profile manually through the [Check Point Infinity Portal](https://portal.checkpoint.com).

Once created, navigate to the profile in the UI and copy the profile ID from the URL — it is in UUID format (e.g. `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`).

### 2. Define the Asset in Terraform

Create the asset in Terraform and attach the `profileID` obtained from step 1. The key fields to configure are:

- `urls` — the public URL(s) to protect
- `upstream_url` — the backend origin URL
- `profiles` — the SaaS profile ID copied from the UI

**Example:**

```hcl
resource "inext_web_app_practice" "demo_practice" {
  name       = "demo-practice"
  visibility = "Shared"
}

resource "inext_web_app_asset" "saas_asset" {
  name         = "demo-asset"
  profiles     = ["<put-the-appsec-saas-profile-id-here>"]
  urls         = ["<put-your-url-here>"]
  upstream_url = "<put-your-upstream-url-here>"

  practice {
    main_mode = "Learn"
    id        = inext_web_app_practice.demo_practice.id
    sub_practices_modes = {
      IPS    = "Disabled"
      WebBot = "AccordingToPractice"
      Snort  = "Disabled"
    }
  }
}
```

For more information see:
- [WAF SaaS deployment guide](https://waf-doc.inext.checkpoint.com/getting-started/deploy-enforcement-point/waf-as-a-service-waf-saas)
- [Security practices concepts](https://waf-doc.inext.checkpoint.com/concepts/security-practices)

### 3. Apply and Publish/Enforce

After defining the asset, run `terraform apply` and then publish and enforce your changes — either via the CLI tool or the `inext_publish_enforce` Terraform resource. See the [root README](../README.md#publish-and-enforce-your-changes-required) for details.

### 4. Certificate Validation (UI Required)

After the asset is applied, navigate to the profile section in the UI and click on the profile. You will see the URL waiting for certificate validation. This domain validation process must be completed through the UI.

### 5. Ongoing Management

Once the certificate is validated, the asset can be fully managed via Terraform.

---

## Known Technical Gaps

The following features are not yet fully supported or have known limitations:

| Gap | Description |
|-----|-------------|
| **Limited practice types** | Not all practice types that can be attached to an AppSec SaaS asset are supported. For example, auth enforcement is not yet available via Terraform. |
| **Trusted Sources (sourceIP)** | Trusted source `sourceIP` configuration for AppSec SaaS assets cannot be configured in the UI, but it _can_ be configured via Terraform. |
| **No custom certificate support** | Bring-your-own-certificate (BYOC) and certificate switching are not supported via Terraform. Certificates are fully managed by WAF SaaS. |
| **Profile creation** | AppSec SaaS profile creation must be done manually in the UI — it is not supported via this provider. |
