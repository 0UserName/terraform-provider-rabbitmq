package provider

import (
	"fmt"
	"strings"
)

const Delimiter = "@"

func GetIdPattern(numberOfSegments int) string {

	pattern := "%s"

	if numberOfSegments > 1 {

		pattern += Delimiter
	}
	return strings.Trim(strings.Repeat(pattern, numberOfSegments), Delimiter)
}

// Format: scope@limit@alias
func createLimitId(scope string, limit string, alias string) string {

	return fmt.Sprintf(GetIdPattern(3), scope, limit, alias)
}

// Format: name@vhost@arguments
func createQueueId(name string, vhost string, arguments string) string {

	return fmt.Sprintf(GetIdPattern(3), name, vhost, arguments)
}

// Format: name@vhost@arguments
func createExchangeId(name string, vhost string, arguments string) string {

	return fmt.Sprintf(GetIdPattern(3), name, vhost, arguments)
}
