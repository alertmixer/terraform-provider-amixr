package amixr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAmixrAction_Basic(t *testing.T) {
	actionName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))
	integrationID := fmt.Sprintf("test-acc-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAmixrActionConfig(actionName, integrationID),
				ExpectError: regexp.MustCompile(`Integration does not exist`),
			},
		},
	})
}

func testAccDataSourceAmixrActionConfig(actionName, integrationID string) string {
	return fmt.Sprintf(`
data "amixr_action" "test-acc-action" {
	name = "%s"
	integration_id = "%s"
}
`, actionName, integrationID)
}
