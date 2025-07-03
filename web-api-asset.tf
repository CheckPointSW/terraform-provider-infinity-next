# Terraform HCL

# Generate 50 practices
locals {
  practice_count = 50
}

resource "inext_web_api_practice" "test_update_practice" {
  count = local.practice_count
  name  = "test update practice${count.index + 1}"
}

# Generate 50 assets, each referencing its corresponding practice
resource "inext_web_api_asset" "test_update_asset" {
  count        = local.practice_count
  name         = "test update asset${count.index + 1}"
  upstream_url = "https://www.acme.com"
  urls         = ["http://example${count.index + 1}.com", "https://example${count.index + 1}.com"]
  profiles     = ["d0cb3b86-4142-1579-aeaa-7eb64139224e"]

  practice {
    main_mode = "Prevent"
    id        = inext_web_api_practice.test_update_practice[count.index].id
  }
}

