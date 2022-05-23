resource "inext_web_app_practice" "stam" {
	name                          = "stam"
	ips {
		performance_impact    = "LowOrLower"
		severity_level        = "LowOrAbove"
		protections_from_year = "2016"
		high_confidence       = "Detect"
		medium_confidence     = "Detect"
		low_confidence        = "Detect"
	}
	web_attacks {
		minimum_severity = "High"
		advanced_setting {
			csrf_protection      = "Prevent"
			open_redirect        = "Disabled"
			error_disclosure     = "AccordingToPractice"
			body_size            = 1000
			url_size             = 1000
			header_size          = 1000
			max_object_depth     = 1000
			illegal_http_methods = true
		}
	}
	web_bot {
		inject_uris = ["url1", "url2", "url3", "url4"]
		valid_uris  = ["url1", "url2"]
	}
}