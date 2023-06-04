---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_binding"
sidebar_current: "docs-rabbitmq-resource-binding"
description: |-
  Creates and manages a binding on a RabbitMQ server.
---

# rabbitmq\_binding

The ``rabbitmq_binding`` resource creates and manages a binding relationship
between a queue an exchange.

## Example Usage

```hcl
resource "rabbitmq_vhost" "test" {
  name = "test"
}

resource "rabbitmq_permissions" "guest" {
  user  = "guest"
  vhost = "${rabbitmq_vhost.test.name}"

  permissions {
    configure = ".*"
    write     = ".*"
    read      = ".*"
  }
}

resource "rabbitmq_exchange" "test" {
  name  = "test"
  vhost = "${rabbitmq_permissions.guest.vhost}"

  settings {
    type        = "fanout"
    durable     = false
    auto_delete = true
  }
}

resource "rabbitmq_queue" "test" {
  name  = "test"
  vhost = "${rabbitmq_permissions.guest.vhost}"

  settings {
    durable     = true
    auto_delete = false
  }
}

resource "rabbitmq_binding" "test" {
  source           = "${rabbitmq_exchange.test.name}"
  vhost            = "${rabbitmq_vhost.test.name}"
  destination      = "${rabbitmq_queue.test.name}"
  destination_type = "queue"
  routing_key      = "#"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) The source exchange.

* `vhost` - (Required) The vhost to create the resource in.

* `destination` - (Required) The destination queue or exchange.

* `destination_type` - (Required) The type of destination (queue or exchange).

* `routing_key` - (Optional) A routing key for the binding.

* `arguments` - (Optional) Additional key/value arguments for the binding.

~> **NOTE:** The source and destination properties take the names of queues or exchangers as arguments. However, it is
acceptable (and desirable) to use the identifiers of these resources. This will help to correctly track the state of the
binding when the state of the mentioned resources changes. See [GH-3][GH-3] (new) [GH-34][GH-34], [GH-25][GH-25].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `properties_key` - A unique key to refer to the binding.

## Import

Bindings can be imported using the `id` which is composed of
`vhost/source/destination/destination_type/properties_key`. E.g.

```
$ terraform import rabbitmq_binding.test test/test/test/queue/%23
```

[GH-3]: https://github.com/0UserName/terraform-provider-rabbitmq/issues/3

[GH-34]: https://github.com/cyrilgdn/terraform-provider-rabbitmq/issues/34

[GH-25]: https://github.com/cyrilgdn/terraform-provider-rabbitmq/issues/25