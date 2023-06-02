---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_queue"
sidebar_current: "docs-rabbitmq-data-source-queue"
description: |-
Provides a queue data source on a RabbitMQ server.
---

# rabbitmq\_queue

The ``rabbitmq_queue`` data source can be used to get the general attributes of queue.

## Example Usage

### Basic Example

```hcl
data "rabbitmq_vhost" "test" {

  name = "/"
}

data "rabbitmq_queue" "test" {

  vhost = rabbitmq_vhost.test.name
  name  = "test"
}
```

## Argument Reference

The following arguments are supported:

* `vhost` - (Required) The vhost to create the resource in.

* `name` - (Required) The name of the queue.

## Attributes Reference

No further attributes are exported.