package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVhostResource(t *testing.T) {

	rn := "rabbitmq_vhost.r_test"
	dn := "rabbitmq_vhost.d_test"

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccReadTestData("vhost_test.tf"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr(rn, "id", "r_test"),
					resource.TestCheckResourceAttr(dn, "id", "d_test"),

					resource.TestCheckResourceAttr(rn, "name", "r_test"),
					resource.TestCheckResourceAttr(dn, "name", "d_test"),
				),
			},

			// ImportState testing
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"queue_type"},
			},
		},
	})
}
