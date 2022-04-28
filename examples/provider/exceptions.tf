resource "inext_exceptions" "test" {
  name = "inext_exceptions-test1"
  exception {
    match = { # currently matches with "AND" condition between all keys and their values
      hostName         = "www.google.com"
      url              = "/login"
      sourceIdentifier = "1.1.1.1/24"
    }
    action  = "drop" # enum of ["drop", "skip", "accept", "suppressLog"]
    comment = "some comment"
  }
  exception {
    match = {
      hostName         = "www.acme.com"
      url              = "/"
      sourceIdentifier = "1.0.0.0/18"
    }
    action = "skip" # enum of ["drop", "skip", "accept", "suppressLog"]
  }
  exception {
    match = {
      hostName         = "www.checkpoint.com"
      url              = "/"
      sourceIdentifier = "1.0.0.0/18"
    }
    action = "accept" # enum of ["drop", "skip", "accept", "suppressLog"]
  }
  exception {
    match = {
      hostName         = "www.apple.com"
      url              = "/"
      sourceIdentifier = "1.0.0.0/18"
    }
    action = "suppressLog" # enum of ["drop", "skip", "accept", "suppressLog"]
  }
}