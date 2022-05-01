package rabbitmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceVhostConfig_basic = `
data "rabbitmq_vhost" "test" {

  name = "/"
}`

func TestAccDataSourceVhost_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() {

			testAccPreCheck(t)
		},

		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVhostConfig_basic,
				Check:  resource.ComposeTestCheckFunc(resource.TestMatchResourceAttr("data.rabbitmq_vhost.test", "id", regexp.MustCompile("^/$"))),
			},
		},
	})
}
