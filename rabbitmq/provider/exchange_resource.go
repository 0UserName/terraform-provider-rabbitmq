package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ resource.Resource = &ExchangeResource{}
var _ resource.ResourceWithImportState = &ExchangeResource{}

func NewExchangeResource() resource.Resource {

	return &ExchangeResource{}
}

type ExchangeResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type ExchangeSettingsResourceModel struct {
	Type       types.String  `tfsdk:"type"`
	Durable    types.Bool    `tfsdk:"durable"`
	AutoDelete types.Bool    `tfsdk:"auto_delete"`
	Arguments  types.Dynamic `tfsdk:"arguments"`
}

type ExchangeResourceModel struct {
	Id       types.String                  `tfsdk:"id"`
	Vhost    types.String                  `tfsdk:"vhost"`
	Name     types.String                  `tfsdk:"name"`
	Settings ExchangeSettingsResourceModel `tfsdk:"settings"`
}

func (r *ExchangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_exchange"
}

func (r *ExchangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Blocks: map[string]schema.Block{

			"settings": schema.SingleNestedBlock{

				Attributes: map[string]schema.Attribute{

					"type": schema.StringAttribute{

						Required: true,
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

func (r *ExchangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *ExchangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ExchangeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create exchange: %v", data.Name))

	res, err := r.client.DeclareExchange(data.Vhost.ValueString(), data.Name.ValueString(), getExchangeSettings(ctx, &data.Settings))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Exchange creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Exchange creation error", err.Error())

		return
	}

	data.Id = types.StringValue(data.Vhost.ValueString() + data.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExchangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ExchangeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract exchange: %v", data.Id))

	exchange, err := r.client.GetExchange(data.Vhost.ValueString(), data.Name.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Exchange extraction response: %#v", exchange))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: Exchange extraction error", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExchangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// NO-OP
}

func (r *ExchangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ExchangeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete exchange: %#v", data.Id))

	res, err := r.client.DeleteExchange(data.Vhost.ValueString(), data.Name.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Exchange deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: Exchange deletion error", err.Error())
	}
}

func (r *ExchangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func getExchangeSettings(ctx context.Context, settings *ExchangeSettingsResourceModel) rabbithole.ExchangeSettings {

	arguments := make(map[string]any)

	_ = json.Unmarshal([]byte(settings.Arguments.String()), &arguments)

	return rabbithole.ExchangeSettings{Type: settings.Type.ValueString(), Durable: settings.Durable.ValueBool(), AutoDelete: settings.AutoDelete.ValueBool(), Arguments: arguments}
}
