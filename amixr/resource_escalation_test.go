package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAmixrEscalation_basic(t *testing.T) {
	riName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))
	reType := "wait"
	reDuration := "300"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAmixrEscalationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAmixrEscalationConfig(riName, reType, reDuration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAmixrEscalationResourceExists("amixr_escalation.foo"),
					resource.TestCheckResourceAttr("amixr_escalation.foo", "type", "wait"),
				),
			},
		},
	})
}

func testAccCheckAmixrEscalationResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*amixr.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "amixr_escalation" {
			continue
		}

		if _, _, err := client.Escalations.GetEscalation(r.Primary.ID, &amixr.GetEscalationOptions{}); err == nil {
			return fmt.Errorf("Escalation still exists")
		}

	}
	return nil
}

func testAccAmixrEscalationConfig(riName, reType, reDuration string) string {
	return fmt.Sprintf(`
resource "amixr_integration" "foo" {
	name = "%s"
	type = "grafana"
}

resource "amixr_escalation" "foo" {
	route_id = amixr_integration.foo.default_route_id
	type = "%s"
	duration = "%s"
	position = 0
}
`, riName, reType, reDuration)
}

func testAccCheckAmixrEscalationResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Escalation ID is set")
		}

		client := testAccProvider.Meta().(*amixr.Client)

		found, _, err := client.Escalations.GetEscalation(rs.Primary.ID, &amixr.GetEscalationOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Escalation policy not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}
