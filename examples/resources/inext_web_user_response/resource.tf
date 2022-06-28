resource "inext_web_user_response" "web-user-response-blockpage" {
  name               = "web-user-response"
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