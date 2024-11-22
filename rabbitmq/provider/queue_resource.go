package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ resource.Resource = &QueueResource{}
var _ resource.ResourceWithImportState = &QueueResource{}

func NewQueueResource() resource.Resource {

	return &QueueResource{}
}

type QueueResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type QueueSettingsResourceModel struct {
	Type       types.String  `tfsdk:"type"`
	Durable    types.Bool    `tfsdk:"durable"`
	AutoDelete types.Bool    `tfsdk:"auto_delete"`
	Arguments  types.Dynamic `tfsdk:"arguments"`
}

type QueueResourceModel struct {
	Id       types.String               `tfsdk:"id"`
	Vhost    types.String               `tfsdk:"vhost"`
	Name     types.String               `tfsdk:"name"`
	Settings QueueSettingsResourceModel `tfsdk:"settings"`
}

func (r *QueueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_queue"
}

func (r *QueueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Blocks: map[string]schema.Block{

			"settings": schema.SingleNestedBlock{

				Attributes: map[string]schema.Attribute{

					"type": schema.StringAttribute{

						Required: true,

						Validators: []validator.String{

							stringvalidator.OneOf("classic", "quorum", "stream"),
						},
					},

					"durable": schema.BoolAttribute{

						Computed: true,
						Optional: true,

						Default: booldefault.StaticBool(true),
					},

					"auto_delete": schema.BoolAttribute{

						Computed: true,
						Optional: true,

						Default: booldefault.StaticBool(false),
					},

					"arguments": schema.DynamicAttribute{

						Optional: true,
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
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

func (r *QueueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *QueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data QueueResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create queue: %v", data.Name))

	res, err := r.client.DeclareQueue(data.Vhost.ValueString(), data.Name.ValueString(), getQueueSettings(ctx, &data.Settings))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Queue creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Queue creation error", err.Error())

		return
	}

	data.Id = types.StringValue(data.Vhost.ValueString() + data.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data QueueResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract queue: %v", data.Id))

	queue, err := r.client.GetQueue(data.Vhost.ValueString(), data.Name.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Queue extraction response: %#v", queue))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: Queue extraction error", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// NO-OP
}

func (r *QueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data QueueResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete queue: %#v", data.Id))

	res, err := r.client.DeleteQueue(data.Vhost.ValueString(), data.Name.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Queue deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: Queue deletion error", err.Error())
	}
}

func (r *QueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func getQueueSettings(ctx context.Context, settings *QueueSettingsResourceModel) rabbithole.QueueSettings {

	arguments := make(map[string]any)

	_ = json.Unmarshal([]byte(settings.Arguments.String()), &arguments)

	return rabbithole.QueueSettings{Type: settings.Type.ValueString(), Durable: settings.Durable.ValueBool(), AutoDelete: settings.AutoDelete.ValueBool(), Arguments: arguments}
}
