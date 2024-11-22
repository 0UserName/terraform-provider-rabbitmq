# For resource units
resource "rabbitmq_vhost" "r_test" {

  name = "r_test"
}


# For data_source units
data "rabbitmq_vhost" "d_test" {

  name = rabbitmq_vhost.r_test.name
}