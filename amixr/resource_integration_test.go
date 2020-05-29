package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAmixrIntegration_basic(t *testing.T) {
	rName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))
	rType := "grafana"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAmixrIntegrationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAmixrIntegrationConfig(rName, rType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAmixrIntegrationResourceExists("amixr_integration.foo"),
					resource.TestCheckResourceAttr("amixr_integration.foo", "name", rName),
					resource.TestCheckResourceAttr("amixr_integration.foo", "type", rType),
				),
			},
		},
	})
}

func testAccCheckAmixrIntegrationResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*amixr.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "amixr_integration" {
			continue
		}

		if _, _, err := client.Integrations.GetIntegration(r.Primary.ID, &amixr.GetIntegrationOptions{}); err == nil {
			return fmt.Errorf("Integration still exists")
		}

	}
	return nil
}

func testAccAmixrIntegrationConfig(rName, rType string) string {
	return fmt.Sprintf(`
resource "amixr_integration" "foo" {
	name = "%s"
	type = "%s"
}
`, rName, rType)
}

func testAccCheckAmixrIntegrationResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Integration ID is set")
		}

		client := testAccProvider.Meta().(*amixr.Client)

		found, _, err := client.Integrations.GetIntegration(rs.Primary.ID, &amixr.GetIntegrationOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Integration not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}
