---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_user"
sidebar_current: "docs-rabbitmq-data-source-user"
description: |-
Provides a user data source on a RabbitMQ server.
---

# rabbitmq\_user

The ``rabbitmq_user`` data source can be used to get the general attributes of user.

## Example Usage

### Basic Example

```hcl
data "rabbitmq_user" "test" {

  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user.

## Attributes Reference

No further attributes are exported.