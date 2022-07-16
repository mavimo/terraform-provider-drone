[![registry Actions Status](https://github.com/mavimo/terraform-provider-drone/workflows/registry/badge.svg)](https://github.com/mavimo/terraform-provider-drone/actions)
# Drone Terraform Provider

A [Terraform](https://www.terraform.io) provider for configuring the [Drone](https://drone.io) continuous delivery platform.


This is a fork of [this Drone provider](https://github.com/KazanExpress/terraform-provider-drone/actions) that is a fork of [this other provider](https://github.com/Lucretius/terraform-provider-drone/releases/latest) which does not appear to be maintained anymore.

## Installing

You can download the plugin from the [Releases](https://github.com/mavimo/terraform-provider-drone/releases/latest) page,
for help installing please refer to the [Official Documentation](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).

## Documentation

You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/mavimo/drone/latest/docs).

## Source

To install from source:

```shell
git clone git://github.com/mavimo/terraform-provider-drone.git
cd terraform-provider-drone
go get
go build
```

## Local development environment

In order to have a fully functional development environment with an instance of drone you can bootstrap all the requirement parts using [docker compose](https://docs.docker.com/compose/).

This command will start all the required elements
```
docker compose -f .github/workflows/acceptance-tests/docker-compose.yaml up -d 
```

than -once all the containers are up & running- execute:
```
docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec gitea-mysql mysql -uroot -pgitea gitea -e 'SOURCE /tmp/gitea-data.sql'
curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"hook-test"}'
docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec drone-mysql mysql -uroot -pdrone drone -e 'SOURCE /tmp/drone-data.sql'
curl -XPOST -H 'Authorization: Bearer 5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL' http://localhost:8000/api/user/repos
```

To have a fully functional instance of drone (running on http://localhost:8000) connected to a gitea (hostname: `http://localhost:3000`, username: `terraform`, password: `terraform`) instance.

Than you can use the [drone cli](https://github.com/harness/drone-cli) using the token `5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL` like:
```
drone -t 5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL -s http://localhost:8000/ info
```

> NOTE: during the drone login you will be redirected to the http://gitea-server:3000/ replace the host name with localhost to continue.

## Licence

This project is licensed under the [MIT licence](http://dan.mit-license.org/).

## Meta

This project uses [Semantic Versioning](http://semver.org/).
