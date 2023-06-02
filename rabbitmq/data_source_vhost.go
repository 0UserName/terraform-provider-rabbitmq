package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func dataSourceVhost() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceVhostRead,

		Schema: map[string]*schema.Schema{

			"name": {

				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceVhostRead(d *schema.ResourceData, meta interface{}) error {

	rmqc := meta.(*rabbithole.Client)

	name, _, _, err := parseIdWithArgs(d.Get("name").(string))

	if err != nil {

		return err
	}

	vhost, err := rmqc.GetVhost(name)

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate vhost: %s", err))
	}

	d.SetId(vhost.Name)

	return nil
}
