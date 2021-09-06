package amixr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAmixrUser_Basic(t *testing.T) {
	uEmail := os.Getenv("AMIXR_TEST_USER_EMAIL")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAmixrUserConfig(uEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.amixr_user.test-acc-user", "email", uEmail),
				),
			},
		},
	})
}

func testAccDataSourceAmixrUserConfig(email string) string {
	return fmt.Sprintf(`
data "amixr_user" "test-acc-user" {
	email = "%s"
}
`, email)
}
