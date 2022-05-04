package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func dataSourceQueue() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceQueueRead,

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

func dataSourceQueueRead(d *schema.ResourceData, meta interface{}) error {

	rmqc := meta.(*rabbithole.Client)

	queue, err := rmqc.GetQueue(d.Get("vhost").(string), d.Get("name").(string))

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate queue: %s", err))
	}

	guid, err := generateUUID()

	if err != nil {

		return err
	}

	d.SetId(fmt.Sprintf("%s@%s@%s", queue.Name, queue.Vhost, guid))

	return nil
}
