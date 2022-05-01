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

	exchange, err := rmqc.GetExchange(d.Get("vhost").(string), d.Get("name").(string))

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate exchange: %s", err))
	}

	guid, err := generateUUID()

	if err != nil {

		return err
	}

	d.SetId(fmt.Sprintf("%s@%s@%s", exchange.Name, exchange.Vhost, guid))

	return nil
}
