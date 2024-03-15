terraform {
  required_providers {
    ccoe-naming = {
      version = "0.0.1"
      source  = "gms.corp/dev/ccoe-naming"
    }
  }
}
provider "ccoe-naming" {}
# data "ccoe-naming_resources" "test" {}
resource "ccoe-naming_resources" "test" {
  product     = "Azure Service Bus"
  function    = "generic"
  application = "Platform"
  region      = "eastus2"
  env         = "production"
}
resource "ccoe-naming_vms" "test" {
  region  = "eastus2"
  env     = "production"
  product = "genericvm"
  os      = "windows"
}
output "test-resource-ccoe-naming" {
  value = ccoe-naming_resources.test
}
output "test-vm-ccoe-naming" {
  value = ccoe-naming_vms.test
}


