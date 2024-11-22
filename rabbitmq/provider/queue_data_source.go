package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QueueDataSource{}

func NewQueueDataSource() datasource.DataSource {
	return &QueueDataSource{}
}

type QueueDataSource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type QueueDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Vhost types.String `tfsdk:"vhost"`
	Name  types.String `tfsdk:"name"`
}

func (d *QueueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_queue"
}

func (d *QueueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,
			},

			"vhost": schema.StringAttribute{

				Required: true,
			},

			"name": schema.StringAttribute{

				Required: true,
			},
		},
	}
}

func (d *QueueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	d.client = client
}

func (d *QueueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data QueueDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	queue, err := d.client.GetQueue(data.Vhost.ValueString(), data.Name.ValueString())

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Queue extraction error", err.Error())

		return
	}

	data.Id = types.StringValue(createExchangeId(queue.Name, queue.Vhost, "TODO: ADD ARGS"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
