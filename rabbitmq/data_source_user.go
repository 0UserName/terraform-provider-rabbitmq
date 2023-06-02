package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func dataSourceUser() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{

			"name": {

				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {

	rmqc := meta.(*rabbithole.Client)

	user, err := rmqc.GetUser(d.Get("name").(string))

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate user: %s", err))
	}

	d.SetId(user.Name)

	return nil
}
