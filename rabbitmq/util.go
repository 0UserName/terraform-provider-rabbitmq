package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func checkDeleted(d *schema.ResourceData, err error) error {
	var errorResponse rabbithole.ErrorResponse
	if errors.As(err, &errorResponse) {
		if errorResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}
	return err
}

// Because slashes are used to separate different components when constructing binding IDs,
// we need a way to ensure any components that include slashes can survive the round trip.
// Percent-encoding is a straightforward way of doing so.
// (reference: https://developer.mozilla.org/en-US/docs/Glossary/percent-encoding)

func percentEncodeSlashes(s string) string {
	// Encode any percent signs, then encode any forward slashes.
	return strings.Replace(strings.Replace(s, "%", "%25", -1), "/", "%2F", -1)
}

func percentDecodeSlashes(s string) string {
	// Decode any forward slashes, then decode any percent signs.
	return strings.Replace(strings.Replace(s, "%2F", "/", -1), "%25", "%", -1)
}

// Parses name, vhost from standard resource id. args
// is an array containing the rest of the id segments.
// NOTE: don't use for binding resource.
func parseIdWithArgs(resourceId string) (name string, vhost string, args []string, err error) {

	parts := strings.Split(resourceId, "@")

	switch len(parts) {

	case 0:
		return "", "", nil, fmt.Errorf("unable to parse resource id: %s", resourceId)
	case 1:
		return parts[0], "", []string{}, nil
	case 2:
		return parts[0], parts[1], []string{}, nil
	default:
		return parts[0], parts[1], parts[2:], nil
	}
}

// Parses the resource name from the standard resource id.
func parseName(id string) string {

	name := strings.Split(id, "@")[0]

	log.Printf("[DEBUG] RabbitMQ: read resource name %s from id: %s", name, id)

	return name
}

// Convert dictionary to JSON string.
func toString(arguments map[string]interface{}) (result string) {

	raw, _ := json.Marshal(arguments)

	return string(raw)
}
