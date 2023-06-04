<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for RabbitMQ

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/0UserName/terraform-provider-rabbitmq?label=release&style=for-the-badge)](https://github.com/0UserName/terraform-provider-rabbitmq/releases/latest) [![License](https://img.shields.io/github/license/0UserName/terraform-provider-rabbitmq.svg?style=for-the-badge)](LICENSE)

The Terraform Provider for RabbitMQ is a plugin for Terraform that allows you to interact with RabbitMQ. This provider can be used to manage virtual hosts, users, permissions, policies, limits, queues, exchanges, bindings, and more.

Learn more:

* Read the provider [documentation][provider-documentation].

## Requirements

* [Terraform 1.1.7+][terraform-install] (may work with previous versions)

  For general information about Terraform, visit [terraform.io][terraform-install] and [the project][terraform-github] on GitHub.

* [Go 1.17][golang-install]

  Required if building the provider.

* [RabbitMQ][rabbitmq-releases] 3.9.x or later.

  The provider supports versions in accordance with the RabbitMQ support policies.

    * Learn more: [RabbitMQ Support Policy][rabbitmq-support-policy]

## Using the Provider

To use a released version of the Terraform provider in your environment, run `terraform init` and Terraform will automatically install the provider from the Terraform Registry. For either installation method, documentation about the provider configuration, resources, and data sources can be found on the Terraform Registry.

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider:

```shell
terraform init -upgrade
```

## License

The Terraform Provider for RabbitMQ is available under the [Mozilla Public License, version 2.0][provider-license] license.

[provider-documentation]: https://registry.terraform.io/providers/0UserName/rabbitmq/latest/docs

[terraform-install]: https://www.terraform.io/downloads.html

[terraform-github]: https://github.com/hashicorp/terraform

[golang-install]: https://golang.org/doc/install

[rabbitmq-support-policy]: https://www.rabbitmq.com/versions.html

[rabbitmq-releases]: https://github.com/rabbitmq/rabbitmq-server/releases

[provider-license]: LICENSE