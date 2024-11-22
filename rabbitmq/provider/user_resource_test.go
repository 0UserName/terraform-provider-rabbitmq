package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {

	rn := "rabbitmq_user.r_test"
	dn := "rabbitmq_user.d_test"

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccReadTestData("user_test.tf"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr(rn, "id", "r_test"),
					resource.TestCheckResourceAttr(dn, "id", "r_test"),

					resource.TestCheckResourceAttr(rn, "name", "r_test"),
					resource.TestCheckResourceAttr(dn, "name", "r_test"),

					resource.TestCheckResourceAttr(rn, "tags", "r_test"),
				),
			},

			// ImportState testing
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
