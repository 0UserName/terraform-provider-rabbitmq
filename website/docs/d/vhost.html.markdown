---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_vhost"
sidebar_current: "docs-rabbitmq-data-source-vhost"
description: |-
Provides a vhost data source on a RabbitMQ server.
---

# rabbitmq\_vhost

The ``rabbitmq_vhost`` data source can be used to get the general attributes of vhost.

## Example Usage

### Basic Example

```hcl
data "rabbitmq_vhost" "test" {

  name = "/"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the queue.

## Attributes Reference

No further attributes are exported.