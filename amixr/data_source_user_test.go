package amixr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr("data.amixr_user.foo", "email", uEmail),
				),
			},
		},
	})
}

func testAccDataSourceAmixrUserConfig(email string) string {
	return fmt.Sprintf(`
data "amixr_user" "foo" {
	email = "%s"
}
`, email)
}
