data "local_file" "schema_validation_file" {
  filename = format("oasschema.yaml")
}

resource "inext_web_api_practice" resource-aaaaac {
		name = "resource-aaaaac"
		ips {
			performance_impact    = "MediumOrLower"
			severity_level        = "MediumOrAbove"
			protections_from_year = "2016"
			high_confidence       = "AccordingToPractice"
			medium_confidence     = "AccordingToPractice"
			low_confidence        = "Detect"
		}
		api_attacks {
			minimum_severity = "High"
			advanced_setting {
				body_size            = 1000000
				url_size             = 32768
				header_size          = 102400
				max_object_depth     = 40
				illegal_http_methods = false
			}
		}
		schema_validation {
			data     = data.local_file.schema_validation_file.content
		}
}
	
resource "inext_web_api_asset" omritheking5 {
	  name         = "omritheking5"
	  urls         = ["https://example2.com"]
	  profiles = ["e6cb0fec-c235-a5ba-e746-06dabca1ab02"]
	  practice {
			main_mode = "Prevent"
			sub_practices_modes = {
					IPS = "Detect"
					WebBot = "Prevent"
					SchemaValidation = "Prevent"
			}
			
			id = inext_web_api_practice.resource-aaaaac.id
	}
	mtls {
		filename = "cert.der"
		certificate_type = ".der"
		data	 = "cert data"
		type = "client"
		enable = true
	}
	additional_instructions_blocks {
		filename = "location.json"
		filename_type = ".json"
		data	 = "location data"
		type = "location_instructions"
		enable = true
	}
	redirect_to_https = "true"
	access_log = "true"
	custom_headers {
		name   = "first header"
		value  = "first value"
	}
}

