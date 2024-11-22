# For resource units
resource "rabbitmq_exchange" "r_test" {

  vhost = "/"
  name  = "r_test"

  settings {

    type = "fanout"

    durable     = true
    auto_delete = false
  }
}


# For data_source units
data "rabbitmq_exchange" "d_test" {

  vhost = rabbitmq_exchange.r_test.vhost
  name  = rabbitmq_exchange.r_test.name
}