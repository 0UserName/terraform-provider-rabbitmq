package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VhostDataSource{}

func NewVhostDataSource() datasource.DataSource {

	return &VhostDataSource{}
}

type VhostDataSource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type VhostDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *VhostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_vhost"
}

func (d *VhostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,
			},

			"name": schema.StringAttribute{

				Required: true,
			},
		},
	}
}

func (d *VhostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	d.client = client
}

func (d *VhostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data VhostDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	vhost, err := d.client.GetVhost(data.Name.ValueString())

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Vhost extraction error", err.Error())

		return
	}

	data.Id = types.StringValue(vhost.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
