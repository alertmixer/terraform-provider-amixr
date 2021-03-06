package amixr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAmixrUserGroup_Basic(t *testing.T) {
	slackHandle := fmt.Sprintf("test-acc-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAmixrUserGroupConfig(slackHandle),
				ExpectError: regexp.MustCompile(`couldn't find a user group`),
			},
		},
	})
}

func testAccDataSourceAmixrUserGroupConfig(slackHandle string) string {
	return fmt.Sprintf(`
data "amixr_user_group" "test-acc-user-group" {
	slack_handle = "%s"
}
`, slackHandle)
}
