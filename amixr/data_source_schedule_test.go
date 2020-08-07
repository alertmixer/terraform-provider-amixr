package amixr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAmixrSchedule_Basic(t *testing.T) {
	scheduleName := fmt.Sprintf("test-acc-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAmixrScheduleConfig(scheduleName),
				ExpectError: regexp.MustCompile(`couldn't find a schedule`),
			},
		},
	})
}

func testAccDataSourceAmixrScheduleConfig(scheduleName string) string {
	return fmt.Sprintf(`
data "amixr_schedule" "test-acc-schedule" {
	name = "%s"
}
`, scheduleName)
}
