package amixr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAmixrTeam_Basic(t *testing.T) {
	teamName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAmixrTeamConfig(teamName),
				ExpectError: regexp.MustCompile(`couldn't find a team`),
			},
		},
	})
}

func testAccDataSourceAmixrTeamConfig(teamName string) string {
	return fmt.Sprintf(`
data "amixr_team" "test-acc-team" {
	name = "%s"
}
`, teamName)
}
