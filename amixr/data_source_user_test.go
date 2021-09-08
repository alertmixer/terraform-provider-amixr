package amixr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAmixrUser_Basic(t *testing.T) {
	username := os.Getenv("AMIXR_TEST_USER_USERNAME")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAmixrUserConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.amixr_user.test-acc-user", "username", username),
				),
			},
		},
	})
}

func testAccDataSourceAmixrUserConfig(username string) string {
	return fmt.Sprintf(`
data "amixr_user" "test-acc-user" {
	username = "%s"
}
`, username)
}
