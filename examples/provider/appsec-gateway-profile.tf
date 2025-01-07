resource "inext_appsec_gateway_profile" "test2" {
  name                          = "inext_appsec_gateway_profile-test1"
  profile_sub_type              = "Aws"        # enum of ["Aws", "Azure", "VMware", "HyperV"]
  upgrade_mode                  = "Scheduled"  # enum of ["Automatic", "Manual", "Scheduled"]
  upgrade_time_schedule_type    = "DaysInWeek" # enum of ["DaysInMonth", "DaysInWeek", "Daily"]
  upgrade_time_hour             = "12:00"
  upgrade_time_duration         = 10
  upgrade_time_week_days        = ["Monday", "Thursday", "Friday"]
  reverseproxy_upstream_timeout = 3600
  reverseproxy_additional_settings = {
    Key3 = "Value3"
    Key4 = "Value4"
  }
  max_number_of_agents = 100
  additional_settings = {
    Key1 = "Value1"
    Key2 = "Value2"
  }
}

resource "inext_appsec_gateway_profile" "test" {
  name                          = "inext_appsec_gateway_profile-test2"
  profile_sub_type              = "Aws"        # enum of ["Aws", "Azure", "VMware", "HyperV"]
  upgrade_mode                  = "Automatic"  # enum of ["Automatic", "Manual", "Scheduled"]
  upgrade_time_schedule_type    = "DaysInWeek" # enum of ["DaysInMonth", "DaysInWeek", "Daily"]
  upgrade_time_hour             = "12:00"
  upgrade_time_duration         = 10
  upgrade_time_week_days        = ["Thursday", "Friday", "Monday"]
  reverseproxy_upstream_timeout = 3600
  reverseproxy_additional_settings = {
    Key3 = "Value5"
    Key4 = "Value4"
  }
  max_number_of_agents = 100
  additional_settings = {
    Key1 = "Value1"
    Key2 = "Value2"
  }
}

resource "inext_appsec_gateway_profile" "test3" {
  name                          = "inext_appsec_gateway_profile-test3"
  profile_sub_type              = "Aws"         # enum of ["Aws", "Azure", "VMware", "HyperV"]
  upgrade_mode                  = "Automatic"   # enum of ["Automatic", "Manual", "Scheduled"]
  upgrade_time_schedule_type    = "DaysInMonth" # enum of ["DaysInMonth", "DaysInWeek", "Daily"]
  upgrade_time_hour             = "12:00"
  upgrade_time_duration         = 10
  upgrade_time_days             = [1, 2, 3, 4, 5, 6, 7]
  reverseproxy_upstream_timeout = 3600
  reverseproxy_additional_settings = {
    Key3 = "Value5"
    Key4 = "Value4"
  }
  max_number_of_agents = 100
  additional_settings = {
    Key1 = "Value1"
    Key2 = "Value2"
  }
}