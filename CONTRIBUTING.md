# Local development

In order to have a fully functional development environment with a local instance of Drone and Gitea you can bootstrap all the required services using [docker compose](https://docs.docker.com/compose/).

This command will start all the required elements
```bash
docker compose -f .github/workflows/acceptance-tests/docker-compose.yaml up -d 
```

than, once all the containers are up & running, execute:
```bash
# Import Gitea data
docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec gitea-mysql mysql -uroot -pgitea gitea -e 'SOURCE /tmp/gitea-data.sql'
# Create repository
curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-1"}'
curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-2"}'
curl -k -X POST "http://terraform:terraform@localhost:3000/api/v1/user/repos" -H "content-type: application/json" --data '{"name":"repo-test-3"}'
# Import Drone data
docker compose -f ".github/workflows/acceptance-tests/docker-compose.yaml" exec drone-mysql mysql -uroot -pdrone drone -e 'SOURCE /tmp/drone-data.sql'
# Force Drone to sync repositories
curl -XPOST -H 'Authorization: Bearer 5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL' http://localhost:8000/api/user/repos
```

You can also use `make create-env` and `make load-fixtures` command as shortcut.

You can access to Drone via `http://localhost:8000` and Gitea (hostname: `http://localhost:3000`, username: `terraform`, password: `terraform`) instance.

Than you can use the [drone cli](https://github.com/harness/drone-cli) with the token `5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL` like:
```bash
drone -t 5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL -s http://localhost:8000/ info
```

> [!WARNING]
> during the drone login you will be redirected to the http://gitea-server:3000/ replace the host name with `localhost` to continue.

## Local test

Once you setup the local environment you can run tests via:
```bash
TF_ACC=1 DRONE_SERVER=http://localhost:8000/ DRONE_TOKEN=5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL DRONE_USER=terraform go test ./...
```