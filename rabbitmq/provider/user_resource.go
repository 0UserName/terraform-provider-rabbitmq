package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {

	return &UserResource{}
}

type UserResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type UserResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	Tags     types.String `tfsdk:"tags"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{

				Required: true,
			},

			"password": schema.StringAttribute{

				Required:  true,
				Sensitive: true,
			},

			"tags": schema.StringAttribute{

				Optional: true,
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create user: %v", data.Name))

	res, err := r.client.PutUser(data.Name.ValueString(), getUserSettings(&data))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: User creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: User creation error", err.Error())

		return
	}

	data.Id = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract user: %v", data.Id))

	user, err := r.client.GetUser(data.Id.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: User extraction response: %#v", user))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: User extraction error", err.Error())
		}
	} else {

		data.Tags = types.StringValue(strings.Join(user.Tags, ", "))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to update user %v", data.Id))

	res, err := r.client.PutUser(data.Name.ValueString(), getUserSettings(&data))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: User update response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: User update error", err.Error())
	}
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete user: %#v", data.Id))

	res, err := r.client.DeleteUser(data.Id.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: user deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: User deletion error", err.Error())
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("password"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("tags"), req, resp)
}

func getUserSettings(data *UserResourceModel) rabbithole.UserSettings {

	return rabbithole.UserSettings{Password: data.Password.ValueString(), Tags: strings.Split(data.Tags.ValueString(), ",")}
}
