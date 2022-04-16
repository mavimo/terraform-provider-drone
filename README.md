[![registry Actions Status](https://github.com/KazanExpress/terraform-provider-drone/workflows/registry/badge.svg)](https://github.com/KazanExpress/terraform-provider-drone/actions)
# Drone Terraform Provider

This is a fork of [this Drone provider](https://github.com/Lucretius/terraform-provider-drone/releases/latest) which does not appear to be maintained anymore.
A [Terraform](https://www.terraform.io) provider for configuring the
[Drone](https://drone.io) continuous delivery platform.

## Installing

You can download the plugin from the [Releases](https://github.com/KazanExpress/terraform-provider-drone/releases/latest) page,
for help installing please refer to the [Official Documentation](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).


## Example

```terraform
provider "drone" {
  server = "https:://ci.example.com/"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoib2N0b2NhdCIsInR5cGUiOiJ1c2VyIn0.Fg0eYxO9x2CfGIvIHDZKhQbCGbRAsSB_iRDJlDEW6vc"
}

resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visibility = "public"
}

resource "drone_secret" "master_password" {
  repository = resource.hello_world.repository
  name       = "master_password"
  value      = "correct horse battery staple"
}

resource "drone_cron" "cron_job" {
  repository = resource.hello_world.repository
  name = "cron_job_1"
  expr = "@monthly"
  branch = "test"
}
```

## Provider

#### Example Usage

```terraform
provider "drone" {
  server = "https://ci.example.com/"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoib2N0b2NhdCIsInR5cGUiOiJ1c2VyIn0.Fg0eYxO9x2CfGIvIHDZKhQbCGbRAsSB_iRDJlDEW6vc"
}
````

#### Argument Reference

* `server` - (Optional) The Drone servers url, It must be provided, but can also
  be sourced from the `DRONE_SERVER` environment variable.
* `token` - (Optional) The Drone servers api token, It must be provided, but can
  also be sourced from the `DRONE_TOKEN` environment variable.

## Resources

### `drone_repo`

Activate and configure a repository.

#### Example Usage

```terraform
resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visibility = "public"
}
```

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `trusted` - (Optional) Repository is trusted (default: `false`).
* `protected` - (Optional) Repository is protected (default: `false`).
* `timeout` - (Optional) Repository timeout (default: `60`).
* `visibility` - (Optional) Repository visibility (default: `private`).
* `configuration` - (Optional) Drone Configuration file (default `.drone.yml`).

### `drone_secret`

Manage a repository secret.

#### Example Usage

```terraform
resource "drone_secret" "master_password" {
  repository = "octocat/hello-world"
  name       = "master_password"
  value      = "correct horse battery staple"
}
````

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `name` - (Required) Secret name.
* `value` - (Required) Secret value.
* `allow_on_pull_request` - (Optional) Allow retrieving the secret on pull requests.  (Default `false`)

### `drone_cron`

Manage a repository cron job.

#### Example Usage

```terraform
resource "drone_cron" "cron_job" {
  repository = "octocat/hello-world"
  name       = "the_cron_job"
  expr       = "@monthly"
  event      = "push"
}
````

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `name` - (Required) Cron job name.
* `event` - (Required) The event for this cron job.  Only allowed value is `push`.
* `expr` - (Required) The cron interval.  Allowed values are `@daily`, `@weekly`, `@monthly`, and `@yearly`.
* `branch` - (Optional) The branch to run this cron job on.  (Default `master`)
* `disabled` - (Optional) Is this cron job disabled?


### `drone_user`

Manage a user.

#### Example Usage

```terraform
resource "drone_user" "octocat" {
  login = "octocat"
}
````

```
Note: In order to use the `drone_user` resource you must have admin privileges within your Drone environment.
```

#### Argument Reference

* `login` - (Required) Login name.
* `admin` - (Required) Is the user an admin?
* `active` - (Required) Is the user active?
* `machine` - (Optional) Is the user a machine?  (Default `false`)

#### Attributes Reference

* `token` - The user's access token

## Source

To install from source:

```shell
git clone git://github.com/KazanExpress/terraform-provider-drone.git
cd terraform-provider-drone
go get
go build
```

## Licence

This project is licensed under the [MIT licence](http://dan.mit-license.org/).

## Meta

This project uses [Semantic Versioning](http://semver.org/).
