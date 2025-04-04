package dns_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/unasra/terraform-provider-nios/internal/utils"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/unasra/nios-go-client/dns"
	"github.com/unasra/terraform-provider-nios/internal/acctest"
)

func TestAccRecordaResource_basic(t *testing.T) {
	var resourceName = "nios_dns_a_record.test"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaBasicConfig(name, "10.0.0.20", "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					// TODO: check and validate these
					// Test Read Only fields
					// Test fields with default value
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_disappears(t *testing.T) {
	resourceName := "nios_dns_a_record.test"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRecordaDestroy(context.Background(), &v),
		Steps: []resource.TestStep{
			{
				Config: testAccRecordaBasicConfig(name, "10.0.0.20", "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					testAccCheckRecordaDisappears(context.Background(), &v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRecordaResource_Comment(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_comment"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaComment(name, "10.0.0.20", "default", "This is a new record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "comment", "This is a new record"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaComment(name, "10.0.0.20", "default", "This is an updated record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "comment", "This is an updated record"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Creator(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_creator"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaCreator(name, "10.0.0.20", "default", "STATIC"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "creator", "STATIC"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaCreator(name, "10.0.0.20", "default", "DYNAMIC"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "creator", "DYNAMIC"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_DdnsPrincipal(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_ddns_principal"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaDdnsPrincipal(name, "10.0.0.20", "default", "DDNS_PRINCIPAL_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ddns_principal", "DDNS_PRINCIPAL_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaDdnsPrincipal(name, "10.0.0.20", "default", "DDNS_PRINCIPAL_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ddns_principal", "DDNS_PRINCIPAL_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_DdnsProtected(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_ddns_protected"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaDdnsProtected(name, "10.0.0.20", "default", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ddns_protected", "false"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaDdnsProtected(name, "10.0.0.20", "default", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ddns_protected", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Disable(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_disable"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaDisable(name, "10.0.0.20", "default", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "disable", "false"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaDisable(name, "10.0.0.20", "default", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "disable", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Extattrs(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_extattrs"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"
	extAttrValue1 := acctest.RandomName()
	extAttrValue2 := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaExtattrs(name, "10.0.0.20", "default", map[string]map[string]string{
					"Site": {
						"value": extAttrValue1,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "extattrs.Site.value", extAttrValue1),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaExtattrs(name, "10.0.0.20", "default", map[string]map[string]string{
					"Site": {
						"value": extAttrValue2,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "extattrs.Site.value", extAttrValue2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_ForbidReclamation(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_forbid_reclamation"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaForbidReclamation(name, "10.0.0.20", "default", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "forbid_reclamation", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaForbidReclamation(name, "10.0.0.20", "default", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "forbid_reclamation", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Ipv4addr(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_ipv4addr"
	var v dns.RecordA
	//name := acctest.RandomName() +  ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaIpv4addr("IPV4ADDR_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ipv4addr", "IPV4ADDR_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaIpv4addr("IPV4ADDR_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ipv4addr", "IPV4ADDR_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Name(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_name"
	var v dns.RecordA
	//name := acctest.RandomName() +  ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaName("NAME_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", "NAME_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaName("NAME_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", "NAME_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_Ttl(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_ttl"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaTtl(name, "10.0.0.20", "default", 10, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ttl", "10"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaTtl(name, "10.0.0.20", "default", 30, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "ttl", "30"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRecordaResource_UseTtl(t *testing.T) {
	var resourceName = "nios_dns_a_record.test_use_ttl"
	var v dns.RecordA
	name := acctest.RandomName() + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRecordaUseTtl(name, "10.0.0.20", "default", "true", 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "use_ttl", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccRecordaUseTtl(name, "10.0.0.20", "default", "false", 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordaExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "use_ttl", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckRecordaExists(ctx context.Context, resourceName string, v *dns.RecordA) resource.TestCheckFunc {
	// Verify the resource exists in the cloud
	var readableAttributes = "aws_rte53_record_info,cloud_info,comment,creation_time,creator,ddns_principal,ddns_protected,disable,discovered_data,dns_name,extattrs,forbid_reclamation,ipv4addr,last_queried,ms_ad_user_data,name,reclaimable,shared_record_group,ttl,use_ttl,view,zone"
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		apiRes, _, err := acctest.NIOSClient.DNSAPI.
			RecordaAPI.
			RecordaReferenceGet(ctx, rs.Primary.ID).
			ReturnFields2(readableAttributes).
			ReturnAsObject(1).
			Execute()
		if err != nil {
			return err
		}
		if !apiRes.HasResult() {
			return fmt.Errorf("expected result to be returned: %s", resourceName)
		}
		*v = apiRes.GetResult()
		return nil
	}
}

func testAccCheckRecordaDestroy(ctx context.Context, v *dns.RecordA) resource.TestCheckFunc {
	// Verify the resource was destroyed
	var readableAttributes = "aws_rte53_record_info,cloud_info,comment,creation_time,creator,ddns_principal,ddns_protected,disable,discovered_data,dns_name,extattrs,forbid_reclamation,ipv4addr,last_queried,ms_ad_user_data,name,reclaimable,shared_record_group,ttl,use_ttl,view,zone"
	return func(state *terraform.State) error {
		_, httpRes, err := acctest.NIOSClient.DNSAPI.
			RecordaAPI.
			RecordaReferenceGet(ctx, utils.ExtractResourceRef(*v.Ref)).
			ReturnAsObject(1).
			ReturnFields2(readableAttributes).
			Execute()
		if err != nil {
			if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
				// resource was deleted
				return nil
			}
			return err
		}
		return errors.New("expected to be deleted")
	}
}

func testAccCheckRecordaDisappears(ctx context.Context, v *dns.RecordA) resource.TestCheckFunc {
	// Delete the resource externally to verify disappears test
	return func(state *terraform.State) error {
		_, err := acctest.NIOSClient.DNSAPI.
			RecordaAPI.
			RecordaReferenceDelete(ctx, utils.ExtractResourceRef(*v.Ref)).
			Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccRecordaBasicConfig(name, ipV4Addr, view string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test" {
	name = %q
	ipv4addr = %q
	view = %q
}
`, name, ipV4Addr, view)
}

func testAccRecordaComment(name, ipV4Addr, view, comment string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_comment" {
	name = %q
	ipv4addr = %q
	view = %q
	comment = %q
}
`, name, ipV4Addr, view, comment)
}

func testAccRecordaCreator(name, ipV4Addr, view, creator string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_creator" {
	name = %q
	ipv4addr = %q
	view = %q  
	creator = %q
}
`, name, ipV4Addr, view, creator)
}

func testAccRecordaDdnsPrincipal(name, ipV4Addr, view, ddnsPrincipal string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_ddns_principal" {
	name = %q
	ipv4addr = %q
	view = %q
	ddns_principal = %q
}
`, name, ipV4Addr, view, ddnsPrincipal)
}

func testAccRecordaDdnsProtected(name, ipV4Addr, view, ddnsProtected string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_ddns_protected" {
    name = %q
	ipv4addr = %q
	view = %q
	ddns_protected = %q
}
`, name, ipV4Addr, view, ddnsProtected)
}

func testAccRecordaDisable(name, ipV4Addr, view, disable string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_disable" {
    name = %q
	ipv4addr = %q
	view = %q
	disable = %q
}
`, name, ipV4Addr, view, disable)
}

func testAccRecordaExtattrs(name, ipV4Addr, view string, extAttrs map[string]map[string]string) string {
	valueStr := ""
	extattrsStr := "{\n"
	for k, v := range extAttrs {
		valueStr += "{\n"
		for k1, v1 := range v {
			valueStr += fmt.Sprintf(`
					%s = %q
		`, k1, v1)
		}
		valueStr += "\t}"
		extattrsStr += fmt.Sprintf(`
			%s = %s
	`, k, valueStr)
	}
	extattrsStr += "\t}"
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_extattrs" {
    name = %q
	ipv4addr = %q
	view = %q
	extattrs = %s
}
`, name, ipV4Addr, view, extattrsStr)
}

func testAccRecordaForbidReclamation(name, ipV4Addr, view, forbidReclamation string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_forbid_reclamation" {
    name = %q
	ipv4addr = %q
	view = %q
	forbid_reclamation = %q
}
`, name, ipV4Addr, view, forbidReclamation)
}

func testAccRecordaIpv4addr(ipv4addr string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_ipv4addr" {
	ipv4addr = %q
}
`, ipv4addr)
}

func testAccRecordaName(name string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_name" {
    name = %q
}
`, name)
}

func testAccRecordaTtl(name, ipV4Addr, view string, ttl int32, use_ttl string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_ttl" {
    name = %q
	ipv4addr = %q
	view = %q
	ttl = %d
	use_ttl = %q
}
`, name, ipV4Addr, view, ttl, use_ttl)
}

func testAccRecordaUseTtl(name, ipV4Addr, view, useTtl string, ttl int32) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_use_ttl" {
    name = %q
	ipv4addr = %q
	view = %q
	use_ttl = %q
	ttl = %d
}
`, name, ipV4Addr, view, useTtl, ttl)
}

func testAccRecordaView(view string) string {
	return fmt.Sprintf(`
resource "nios_dns_a_record" "test_view" {
    view = %q
}
`, view)
}
