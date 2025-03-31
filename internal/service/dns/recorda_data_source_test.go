package dns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/unasra/nios-go-client/dns"
	"github.com/unasra/terraform-provider-nios/internal/acctest"
)

func TestAccRecordaDataSource_Filters(t *testing.T) {
	dataSourceName := "data.nios_dns_a_records.test"
	resourceName := "nios_dns_a_record.test"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRecordaDestroy(context.Background(), &v),
		Steps: []resource.TestStep{
			{
				Config: testAccRecordaDataSourceConfigFilters(name, "10.0.0.20", "default"),
				Check: resource.ComposeTestCheckFunc(
					append([]resource.TestCheckFunc{
						testAccCheckRecordaExists(context.Background(), resourceName, &v),
					}, testAccCheckRecordaResourceAttrPair(resourceName, dataSourceName)...)...,
				),
			},
		},
	})
}

func TestAccRecordaDataSource_TagFilters(t *testing.T) {
	dataSourceName := "data.nios_dns_a_records.test"
	resourceName := "nios_dns_a_record.test"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"
	extAttrValue := acctest.RandomName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRecordaDestroy(context.Background(), &v),
		Steps: []resource.TestStep{
			{
				Config: testAccRecordaDataSourceConfigTagFilters(name, "10.0.0.20", "default", extAttrValue),
				Check: resource.ComposeTestCheckFunc(
					append([]resource.TestCheckFunc{
						testAccCheckRecordaExists(context.Background(), resourceName, &v),
					}, testAccCheckRecordaResourceAttrPair(resourceName, dataSourceName)...)...,
				),
			},
		},
	})
}

// below all TestAcc functions

func testAccCheckRecordaResourceAttrPair(resourceName, dataSourceName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(resourceName, "comment", dataSourceName, "result.0.comment"),
		resource.TestCheckResourceAttrPair(resourceName, "creator", dataSourceName, "result.0.creator"),
		resource.TestCheckResourceAttrPair(resourceName, "ddns_principal", dataSourceName, "result.0.ddns_principal"),
		resource.TestCheckResourceAttrPair(resourceName, "ddns_protected", dataSourceName, "result.0.ddns_protected"),
		resource.TestCheckResourceAttrPair(resourceName, "disable", dataSourceName, "result.0.disable"),
		resource.TestCheckResourceAttrPair(resourceName, "extattrs", dataSourceName, "result.0.extattrs"),
		resource.TestCheckResourceAttrPair(resourceName, "forbid_reclamation", dataSourceName, "result.0.forbid_reclamation"),
		resource.TestCheckResourceAttrPair(resourceName, "ipv4addr", dataSourceName, "result.0.ipv4addr"),
		resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "result.0.name"),
		resource.TestCheckResourceAttrPair(resourceName, "ttl", dataSourceName, "result.0.ttl"),
		resource.TestCheckResourceAttrPair(resourceName, "use_ttl", dataSourceName, "result.0.use_ttl"),
		resource.TestCheckResourceAttrPair(resourceName, "view", dataSourceName, "result.0.view"),
	}
}

func testAccRecordaDataSourceConfigFilters(name, ipV4Addr, view string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test" {
	name = %q
	ipv4addr = %q
	view = %q
}

data "nios_dns_a_records" "test" {
	body = {
		"name": nios_dns_a_record.test.name
	}
}
`, name, ipV4Addr, view)
}

func testAccRecordaDataSourceConfigTagFilters(name, ipV4Addr, view, extAttrsValue string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test" {
	name = %q
	ipv4addr = %q
	view = %q
	extattrs = {
		Site = {
			value = %q
		}
	}
}

data "nios_dns_a_records" "test" {
	body = {
		"*Site" = nios_dns_a_record.test.extattrs.Site.value
	}
}
`, name, ipV4Addr, view, extAttrsValue)
}
