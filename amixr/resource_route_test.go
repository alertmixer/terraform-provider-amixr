package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAmixrRoute_basic(t *testing.T) {
	riName := fmt.Sprintf("integration-%s", acctest.RandString(8))
	rrRegex := fmt.Sprintf("regex-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAmixrRouteResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAmixrRouteConfig(riName, rrRegex),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAmixrRouteResourceExists("amixr_route.test-acc-route"),
				),
			},
		},
	})
}

func testAccCheckAmixrRouteResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*amixr.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "amixr_route" {
			continue
		}

		if _, _, err := client.Routes.GetRoute(r.Primary.ID, &amixr.GetRouteOptions{}); err == nil {
			return fmt.Errorf("Route still exists")
		}

	}
	return nil
}

func testAccAmixrRouteConfig(riName string, rrRegex string) string {
	return fmt.Sprintf(`
resource "amixr_integration" "test-acc-integration" {
	name = "%s"
	type = "grafana"
}

resource "amixr_route" "test-acc-route" {
	integration_id = amixr_integration.test-acc-integration.id
	routing_regex = "%s"
	position = 0
}
`, riName, rrRegex)
}

func testAccCheckAmixrRouteResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route ID is set")
		}

		client := testAccProvider.Meta().(*amixr.Client)

		found, _, err := client.Routes.GetRoute(rs.Primary.ID, &amixr.GetRouteOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Route policy not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}
