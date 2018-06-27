# Prometheus exporter for Docker Hub stats

A small Go application to periodically pull statistics and metrics from Docker Hub, and expose them in a [Prometheus](https://prometheus.io/)-compatible format.

## Usage

You can easily run the application from its [Docker image](https://hub.docker.com/r/rycus86/dockerhub-exporter/):

```shell
$ docker run --rm -it -p 8080:8080 rycus86/dockerhub-exporter <flags>
```

The command line parameters are like this:

```
Usage of /exporter:
  -interval duration
        Interval between checks (default 1h0m0s)
  -owner value
        Owners (namespaces) to list repositories for (multiple values are allowed)
  -port int
        The HTTP port to listen on (default 8080)
  -timeout duration
        HTTP API call timeout (default 15s)
```

The Docker image reference points to a multi-arch manifest, with the actual images being available for the `amd64`, `armhf` and `arm64v8` platforms.

### Selecting targets

You can specify the Docker Hub owners (namespaces) to list repositories from multiple times:

```shell
$ docker run --rm -it -p 8080:8080 rycus86/dockerhub-exporter \
      -owner one -owner two -owner three
```

## Metrics

The following metrics are exposed on the `/metrics` endpoint:

```shell
$ curl -s http://localhost:8080/metrics | grep dockerhub_
# HELP dockerhub_pull_count Number of Pulls
# TYPE dockerhub_pull_count gauge
dockerhub_pull_count{owner="rycus86",repository="podlike"} 2379
dockerhub_pull_count{owner="rycus86",repository="prometheus"} 339913
dockerhub_pull_count{owner="rycus86",repository="prometheus-node-exporter"} 89127
# HELP dockerhub_repo_count Number of Repositories
# TYPE dockerhub_repo_count gauge
dockerhub_repo_count{owner="rycus86"} 38
# HELP dockerhub_star_count Number of Stars
# TYPE dockerhub_star_count gauge
dockerhub_star_count{owner="rycus86",repository="podlike"} 0
dockerhub_star_count{owner="rycus86",repository="prometheus"} 2
dockerhub_star_count{owner="rycus86",repository="prometheus-node-exporter"} 1
```

## License

MIT
