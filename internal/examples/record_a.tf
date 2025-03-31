// Create the Record
resource "nios_dns_a_record" "create_record" {
  name     = "example_test3.example.com"
  ipv4addr = "1.2.1.2"
  view     = "default"
  extattrs = {
    "Site" = {
      "value" = "Siteblr"
    }
  }
}

// Read the Record using External Attributes
data "nios_dns_a_records" "read_record_via_extattrs" {
  body = {
    "*Site" : "Siteblr"
  }
  depends_on = [nios_dns_a_record.create_record]
}

output "read_record_via_extattrs" {
  value = data.nios_dns_a_records.read_record_via_extattrs.result
}

// Read Record By Filtering by Name
data "nios_dns_a_records" "read_all_records" {
    body = {
        "name" : "example_test3.example.com"
    }
}

output "read_record_via_name" {
  value = data.nios_dns_a_records.read_record_via_extattrs.result
}