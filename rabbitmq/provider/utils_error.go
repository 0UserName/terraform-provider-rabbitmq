package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func checkDeleted(ctx context.Context, d *tfsdk.State, err error) error {

	var errorResponse rabbithole.ErrorResponse

	if errors.As(err, &errorResponse) {

		if errorResponse.StatusCode == 404 {

			d.RemoveResource(ctx)

			return nil
		}
	}

	return err
}
