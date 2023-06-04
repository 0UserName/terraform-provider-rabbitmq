---
layout: "rabbitmq"
page_title: "Binding state"
sidebar_current: "docs-rabbitmq-guide-binding-state"
description: |-
  Binding state

---

# Binding state

Since RabbitMQ doesn't support updating resources such as `queue` and `exchange`,
then when you try to update their configuration, the resources are re-created. This causes
the bindings to these resources created earlier to be removed, but terraform won't be able
to track these changes. As the result, the state and the actual configuration will differ.

## Solution

In order to get around this problem, it was decided to change the format of the identifier for these resources:

- `name@vhost@arguments` for queue and exchange
- `vhost@destination_type@properties_key@source_id@destination_id` for binding

Thus, if a resource identifier (rather than its name) is used as
the value of the `source` and/or `destination` parameter,
then changing one of them will cause the binding to be recreated:

```json
{
   "attributes":{
      "arguments":null,
      "arguments_json":null,
      "destination":"QueueName@VhostName@{\"arguments\":{\"x-queue-type\":\"classic\"},\"auto_delete\":false,\"durable\":true}",
      "destination_type":"queue",
      "id":"VhostName#queue#%23#ExchangeName@VhostName@{\"arguments\":{},\"auto_delete\":false,\"durable\":true,\"type\":\"topic\"}#QueueName@VhostName@{\"arguments\":{\"x-queue-type\":\"classic\"},\"auto_delete\":false,\"durable\":true}",
      "properties_key":"%23",
      "routing_key":"#",
      "source":"ExchangeName@VhostName@{\"arguments\":{},\"auto_delete\":false,\"durable\":true,\"type\":\"topic\"}",
      "vhost":"VhostName"
   }
}
```