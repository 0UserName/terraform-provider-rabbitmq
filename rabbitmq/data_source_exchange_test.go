package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceExchangeConfig_basic = `
resource "rabbitmq_vhost" "test" {

  name = "test"
}

resource "rabbitmq_permissions" "guest" {

  vhost = rabbitmq_vhost.test.name
  user  = "guest"

  permissions {

    configure = ".*"
    write     = ".*"
    read      = ".*"
  }
}

resource "rabbitmq_exchange" "test" {

  vhost = rabbitmq_permissions.guest.vhost
  name  = "test"

  settings {

    type        = "fanout"
    durable     = false
    auto_delete = true
  }
}

data "rabbitmq_exchange" "test" {

  vhost = rabbitmq_vhost.test.name
  name  = rabbitmq_exchange.test.name
}`

func TestAccDataSourceExchange_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() {

			testAccPreCheck(t)
		},

		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceExchangeConfig_basic,
				Check:  resource.TestMatchResourceAttr("data.rabbitmq_exchange.test", "id", regexp.MustCompile(fmt.Sprintf("^%s@%s@%s$", "test", "test", UuidRegex))),
			},
		},
	})
}
