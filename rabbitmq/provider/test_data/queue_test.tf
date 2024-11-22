# For resource units
resource "rabbitmq_queue" "r_test" {

  vhost = "/"
  name  = "r_test"

  settings {

    type = "quorum"

    durable     = true
    auto_delete = false

   arguments = {

      x-message-ttl: 5000
    }
  }
}


# For data_source units
#data "rabbitmq_queue" "d_test" {

#  vhost = rabbitmq_queue.r_test.vhost
#  name  = rabbitmq_queue.r_test.name
#}