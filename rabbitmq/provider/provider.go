package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ provider.Provider = &RabbitMqProvider{}

type RabbitMqProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type RabbitMqProviderModel struct {
	Endpoint       types.String `tfsdk:"endpoint"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Insecure       types.Bool   `tfsdk:"insecure"`
	CACertFile     types.String `tfsdk:"cacert_file"`
	ClientCertFile types.String `tfsdk:"clientcert_file"`
	ClientKeyFile  types.String `tfsdk:"clientkey_file"`
	Proxy          types.String `tfsdk:"proxy"`
}

func (p *RabbitMqProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {

	resp.TypeName = "rabbitmq"

	resp.Version = p.version
}

func (p *RabbitMqProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"endpoint": schema.StringAttribute{

				Optional: true,
			},

			"username": schema.StringAttribute{

				Optional: true,
			},

			"password": schema.StringAttribute{

				Optional:  true,
				Sensitive: true,
			},

			"insecure": schema.BoolAttribute{

				Optional: true,
			},

			"cacert_file": schema.StringAttribute{

				Optional: true,
			},

			"clientcert_file": schema.StringAttribute{

				Optional: true,
			},

			"clientkey_file": schema.StringAttribute{

				Optional: true,
			},

			"proxy": schema.StringAttribute{

				Optional: true,
			},
		},
	}
}

func (p *RabbitMqProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data RabbitMqProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	if data.Endpoint.IsNull() {

		data.Endpoint = types.StringValue(os.Getenv("RABBITMQ_ENDPOINT"))
	}

	if data.Username.IsNull() {

		data.Username = types.StringValue(os.Getenv("RABBITMQ_USERNAME"))
	}

	if data.Password.IsNull() {

		data.Password = types.StringValue(os.Getenv("RABBITMQ_PASSWORD"))
	}

	if data.Endpoint.IsNull() || data.Username.IsNull() || data.Password.IsUnknown() {

		panic(" ENDPOINT, USERNAME and PASSWORD must be set")
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	//

	// Configure TLS/SSL:
	// Ignore self-signed cert warnings
	// Specify a custom CA / intermediary cert
	// Specify a certificate and key
	tlsConfig := &tls.Config{}

	// CACertFile
	if !data.CACertFile.IsNull() {

		caCert, err := os.ReadFile(data.CACertFile.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Certificate authority error", err.Error())

			return
		}

		caCertPool := x509.NewCertPool()

		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	// ClientCert
	if !data.ClientCertFile.IsNull() && data.ClientKeyFile.IsNull() {

		clientPair, err := tls.LoadX509KeyPair(data.ClientCertFile.ValueString(), data.ClientKeyFile.ValueString())

		if err != nil {

			resp.Diagnostics.AddError("Certificate authority error", err.Error())

			return
		}

		tlsConfig.Certificates = []tls.Certificate{clientPair}
	}

	if !data.Insecure.ValueBool() {

		tlsConfig.InsecureSkipVerify = true
	}

	var proxyURL *url.URL
	if !data.Proxy.IsNull() {

		var err error

		proxyURL, err = url.Parse(data.Proxy.ValueString())

		if err != nil {

			resp.Diagnostics.AddError("Proxy URL parse error", err.Error())

			return
		}
	}

	// Connect to RabbitMQ management interface
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy: func(req *http.Request) (*url.URL, error) {
			if proxyURL != nil {
				return proxyURL, nil
			}

			return http.ProxyFromEnvironment(req)
		},
	}

	client, err := rabbithole.NewTLSClient(data.Endpoint.ValueString(), data.Username.ValueString(), data.Password.ValueString(), transport)

	if err != nil {

		resp.Diagnostics.AddError("Connection error", err.Error())

		if resp.Diagnostics.HasError() {

			return
		}
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RabbitMqProvider) Resources(ctx context.Context) []func() resource.Resource {

	return []func() resource.Resource{

		NewVhostResource,
		NewUserResource,
		NewExchangeResource,
		NewQueueResource,
	}
}

func (p *RabbitMqProvider) DataSources(ctx context.Context) []func() datasource.DataSource {

	return []func() datasource.DataSource{

		NewVhostDataSource,
		NewUserDataSource,
		NewExchangeDataSource,
		NewQueueDataSource,
	}
}

func New(version string) func() provider.Provider {

	return func() provider.Provider {

		return &RabbitMqProvider{

			version: version,
		}
	}
}
