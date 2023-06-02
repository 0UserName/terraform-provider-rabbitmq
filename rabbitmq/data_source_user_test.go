package rabbitmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceUserConfig_basic = `
data "rabbitmq_user" "test" {

  name = "guest"
}`

func TestAccDataSourceUser_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() {

			testAccPreCheck(t)
		},

		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserConfig_basic,
				Check:  resource.ComposeTestCheckFunc(resource.TestMatchResourceAttr("data.rabbitmq_user.test", "id", regexp.MustCompile("guest"))),
			},
		},
	})
}
