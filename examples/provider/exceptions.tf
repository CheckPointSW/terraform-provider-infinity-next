resource "inext_exceptions" "test" {
  name = "inext_exceptions_test1"
  exception {
    match {
      operator = "or" # enum of ["and", "or", "not-equals", "equals", "in", "not-in", "exist"]
      operand {
        operator = "equals"
        key      = "hostName" # enum of ["hostName", "sourceIdentifier", "url", "countryCode", "countryName", "manufacturer", "paramName", "paramValue", "protectionName", "sourceIP"]
        value    = ["www.google.com"]
      }
      operand {
        operator = "and"
        operand {
          operator = "in"
          key      = "sourceIdentifier"
          value    = ["1.1.1.1/24"]
        }
        operand {
          operator = "not-in"
          key      = "countryName"
          value    = ["Ukraine", "Russia"]
        }
      }
      operand {
        key   = "url"
        value = ["/login"]
      }
    }
    action  = "skip" # enum of ["drop", "skip", "accept", "suppressLog"]
    comment = "some comment"
  }
}