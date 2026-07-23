# Deploying Dragonfly with Container Directly

> Currently, docker compose deploying is tested just in single host, no HA support.

## Deploy with Docker Compose

The `run.sh` script will generate config and deploy all components with `docker-compose`.

Just run:

```shell
# Without network=host mode,the HOST IP would be the docker network gateway IP address, use the command below to
# obtain the ip address, "docker network inspect bridge -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}'"
export IP=<host ip>
./run.sh
```

## Configure with environment variables

In addition to the config files under `config/`, the `manager` and `scheduler`
services can have any config value overridden with an environment variable. This
is useful for tweaking a single setting without editing the generated config.

The environment variable name is built from the config key as follows:

- Prefix it with the service name (`MANAGER_` or `SCHEDULER_`).
- Join nested keys with `_` and uppercase the whole thing.
- camelCase keys collapse to all-lowercase (no separating underscore), e.g.
  `advertiseIP` becomes `ADVERTISEIP`.

In other words, the YAML path `a.b.cDef` maps to `<SERVICE>_A_B_CDEF`.

Examples:

| Service   | Config key             | Environment variable            |
| --------- | ---------------------- | ------------------------------- |
| scheduler | `server.port`          | `SCHEDULER_SERVER_PORT`         |
| scheduler | `server.advertiseIP`   | `SCHEDULER_SERVER_ADVERTISEIP`  |
| manager   | `server.rest.addr`     | `MANAGER_SERVER_REST_ADDR`      |
| manager   | `server.grpc.port.start` | `MANAGER_SERVER_GRPC_PORT_START` |

For example, to override the scheduler listen port:

```shell
SCHEDULER_SERVER_PORT=8003 ./run.sh
```

> Note: this applies to the `manager` and `scheduler` services only. The
> `client` and `seed-client` services are configured through their config files.

## Delete containers with docker compose

```shell
docker compose down
```
