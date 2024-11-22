package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueueResource(t *testing.T) {

	rn := "rabbitmq_queue.r_test"

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccReadTestData("queue_test.tf"),

				Check: resource.ComposeAggregateTestCheckFunc(

					//	resource.TestCheckResourceAttr(rn, "id", "r_test"),
					resource.TestCheckResourceAttr(rn, "vhost", "/"),
					resource.TestCheckResourceAttr(rn, "name", "r_test"),
					resource.TestCheckResourceAttr(rn, "settings.type", "quorum"),
					resource.TestCheckResourceAttr(rn, "settings.durable", "true"),
					resource.TestCheckResourceAttr(rn, "settings.auto_delete", "false"),
					resource.TestCheckResourceAttr(rn, "settings.arguments.x-message-ttl", "5000"),
				),
			},

			// ImportState testing
			/*	{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,

				//ImportStateVerifyIgnore: []string{"queue_type"},
			},*/
		},
	})
}
