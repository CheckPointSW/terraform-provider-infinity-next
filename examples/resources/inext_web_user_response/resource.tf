terraform {
  required_providers {
    inext = {
      source = "CheckPointSW/infinity-next"
      version = "1.1.1"
    }
  }
}

provider "inext" {
  region = "eu"
  # client_id  = ""  // can be set with env var INEXT_CLIENT_ID
  # access_key = "" // can be set with env var INEXT_ACCESS_KEY
}

resource "inext_web_user_response" "web-user-response-blockpage" {
  name               = "web-user-response"
  visibility         = "Shared" # enum of ["Shared", "Local"]
  mode               = "BlockPage"
  http_response_code = 403
  message_title      = "some message title"
  message_body       = "some message body"
}

resource "inext_web_user_response" "web-user-response-redirect" {
  name         = "web-user-response-redirect"
  mode         = "Redirect"
  redirect_url = "http://localhost:1234/test"
  x_event_id   = true
}

resource "inext_web_user_response" "web-user-response-responsecodeonly" {
  name               = "web-user-response-responsecodeonly"
  mode               = "ResponseCodeOnly"
  http_response_code = 403
}