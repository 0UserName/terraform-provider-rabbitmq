---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_exchange"
sidebar_current: "docs-rabbitmq-data-source-exchange"
description: |-
Provides an exchange data source on a RabbitMQ server.
---

# rabbitmq\_exchange

The ``rabbitmq_exchange`` data source can be used to get the general attributes of exchange.

## Example Usage

### Basic Example

```hcl
data "rabbitmq_vhost" "test" {

  name = "/"
}

data "rabbitmq_exchange" "test" {

  vhost = "rabbitmq_vhost.test.name"
  name  = "test"
}
```

## Argument Reference

The following arguments are supported:

* `vhost` - (Required) The vhost to create the resource in.

* `name` - (Required) The name of the exchange.

## Attributes Reference

No further attributes are exported.