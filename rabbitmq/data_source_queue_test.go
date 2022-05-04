package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceQueueConfig_basic = `
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

resource "rabbitmq_queue" "test" {

  vhost = rabbitmq_permissions.guest.vhost
  name  = "test"

  settings {

    durable     = false
    auto_delete = true
  }
}

data "rabbitmq_queue" "test" {

  vhost = rabbitmq_vhost.test.name
  name  = rabbitmq_queue.test.name
}`

func TestAccDataSourceQueue_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() {

			testAccPreCheck(t)
		},

		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQueueConfig_basic,
				Check:  resource.TestMatchResourceAttr("data.rabbitmq_queue.test", "id", regexp.MustCompile(fmt.Sprintf("^%s@%s@%s$", "test", "test", UuidRegex))),
			},
		},
	})
}
