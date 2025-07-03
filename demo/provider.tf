terraform {
  required_providers {
    inext = {
      source  = "CheckPointSW/infinity-next"
      version = "1.1.5"
    }
  }
}

provider "inext" {
  region = "dev"
  client_id  = "f34fd2cbfe9341289a369965db3d630b" 
  access_key = "88b86823b92a4a7e9f57c94667a57866"
}