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

	vhost, _, _, err := parseIdWithArgs(d.Get("vhost").(string))

	if err != nil {

		return err
	}

	name, _, _, err := parseIdWithArgs(d.Get("name").(string))

	if err != nil {

		return err
	}

	queue, err := rmqc.GetQueue(vhost, name)

	if err != nil {

		return checkDeleted(d, fmt.Errorf("cannot locate queue: %s", err))
	}

	d.SetId(fmt.Sprintf("%s@%s@%s", queue.Name, queue.Vhost, fmt.Sprintf("%t:%t:%s", queue.Durable, queue.AutoDelete, toString(queue.Arguments))))

	return nil
}
