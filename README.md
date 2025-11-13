# Check Point Infinity Next Management Terraform Provider

Infinity Next's Terraform Provider for managing CloudGuard AppSec and other Infinity Next security application using Terraform.
You could read the documentation of Infinity Next [here](https://appsec-doc.inext.checkpoint.com/).

## Requirements

- Terraform v0.13+
- [inext CLI](https://github.com/CheckPointSW/infinity-next-terraform-cli/releases/latest) - used to publish and enforce your changes made by Terraform.

## Usage

### Generating an API Key _(Required)_

1. Go to https://portal.checkpoint.com, navigate to _Global Settings -> API Keys_

2. Create a new API key and select _Infinity Policy_ as the service, with _Admin_ role, we recommend that you specify a meaningful comment for the key so you could identify them later and avoid mistakes.

3. Store the _Client ID_ and _Secret Key_ in a secure location, and note there's no way to view the secret key afterward.

### Configuring the Provider

There are 2 options to set the credentials to be used by the provider:

- Set the credentials in environment variables `INEXT_CLIENT_ID` and `INEXT_ACCESS_KEY`

- Set the credentials explicitly or through input variables, in the `.tf` file that defines the `provider` block using the fields `client_id` and `access_key`

Note that credentials are per region, which can be configured with the `region` field of the provider's definition. It defaults to "eu" and currently it accepts either "eu" or "us".

### Publish and Enforce your changes _(Required)_

All changes that are made when running `terraform apply` are done under a session of the configured API key.
At Infinity Next, each session must be published to be able to enforce your configured policies on your assets. Think of it as commiting your changes to be able to make a release.

Due to Terraform's lack of concept of session management/commiting changes at the end of an applied configuration, it's required from the user of this provider to publish and enforce the applied configuration by himself.

This repository includes a CLI utility for this exact use case, which includes 2 commands: `publish` and `enforce`.

### Using the `inext` CLI

Download and install the CLI found in the [latest release](https://github.com/CheckPointSW/infinity-next-terraform-cli/releases/latest)

You could run `inext help` and get all available options and commands.

The CLI requires the same credentials used to configure the provider, there are 3 options to pass these credentials to the CLI:

1. Set the environment variables: `INEXT_REGION`, `INEXT_CLIENT_ID` and `INEXT_ACCESS_KEY` and run `inext <command>`, this is more comfortable for usage right after `terraform apply` since it uses the same environment variables.
2. Set credentials using flags `--client-id` (shorthand `-c`) and `--access-key` (shorthand `-k`)

   ```
   inext publish -c $INEXT_CLIENT_ID -k $INEXT_ACCESS_KEY -r us
   ```

3. Create a yaml file at `~/.inext.yaml` with the following content:
   ```
   client-id: <INEXT_CLIENT_ID>
   access-key: <INEXT_ACCESS_KEY>
   region: eu
   ```
   Run `inext <command>` and the CLI would be configured using `~/.inext.yaml` by default, can be set using `inext --config <config-path> <command>`

## Example

```
terraform init
terraform apply
inext publish && inext enforce
```

Then navigate to the [Cloud Tab](https://portal.checkpoint.com/dashboard/policy#/cloud/getting-started) and enable the _Tech Preview_ toggle at the bottom right.
You should now see your applied objects, each in its own tab.

## Build

### Requirements

- Go 1.24+

To build the provider run:

```
go build
```

To build the build the CLI run:

```
cd cmd
go build -o inext
cp inext /usr/local/bin
```
