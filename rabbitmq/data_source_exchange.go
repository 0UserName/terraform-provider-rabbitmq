package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func dataSourceExchange() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceExchangeRead,

		Schema: map[string]*schema.Schema{

			"name": {

				Type:     schema.TypeString,
				Required: true,
			},

			"vhost": {

				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceExchangeRead(d *schema.ResourceData, meta interface{}) error {

	rmqc := meta.(*rabbithole.Client)

	vhost, _, _, err := parseIdWithArgs(d.Get("vhost").(string))

	if err != nil {

		return err
	}

	name, _, _, err := parseIdWithArgs(d.Get("name").(string))

	if err != nil {

		return err
	}

	exchange, err := rmqc.GetExchange(vhost, name)

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate exchange: %s", err))
	}

	d.SetId(fmt.Sprintf("%s@%s@%s", exchange.Name, exchange.Vhost, fmt.Sprintf("%t:%t:%s", exchange.Durable, exchange.AutoDelete, toString(exchange.Arguments))))

	return nil
}
