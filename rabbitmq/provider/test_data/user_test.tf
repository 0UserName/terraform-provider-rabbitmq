# For resource units
resource "rabbitmq_user" "r_test" {

  name     = "r_test"
  password = "r_test"

  tags = "r_test"
}


# For data_source units
#data "rabbitmq_user" "d_test" {

#  name = rabbitmq_user.r_test.name
#}