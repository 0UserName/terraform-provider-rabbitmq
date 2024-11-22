// To run these acceptance tests, you will need access to a RabbitMQ server
// with the management plugin enabled.
//
// Set the RABBITMQ_ENDPOINT, RABBITMQ_USERNAME, and RABBITMQ_PASSWORD
// environment variables before running the tests.
//
// You can run the tests like this:
//    make testacc TEST=./builtin/providers/rabbitmq

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){

	"rabbitmq": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {

	for _, name := range []string{"RABBITMQ_ENDPOINT", "RABBITMQ_USERNAME", "RABBITMQ_PASSWORD"} {

		if v := os.Getenv(name); v == "" {

			t.Fatal("RABBITMQ_ENDPOINT, RABBITMQ_USERNAME and RABBITMQ_PASSWORD must be set for acceptance tests")
		}
	}
}

func testAccReadTestData(name string) string {

	bytes, _ := os.ReadFile(fmt.Sprintf("./test_data/%v", name))

	return string(bytes)
}
