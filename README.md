# Note: This repo is WIP (not usable yet)
# README

**BEFORE YOU START**, please be aware that there are more ways to integrate with your service that don't require creating a service from this template, see https://keptn.sh/docs/0.10.x/integrations/how_integrate/ for more details.

Examples:

* Webhooks: https://keptn.sh/docs/0.10.x/integrations/webhooks/
* Job-Executor-Service: https://github.com/keptn-sandbox/job-executor-service

---

This is a Keptn Service Template written in GoLang. Follow the instructions below for writing your own Keptn integration.

Quick start:

1. In case you want to contribute your service to keptn-sandbox or keptn-contrib, make sure you have read and understood the [Contributing Guidelines](https://github.com/keptn-sandbox/contributing).
1. Click [Use this template](https://github.com/keptn-sandbox/sumologic-service/generate) on top of the repository, or download the repo as a zip-file, extract it into a new folder named after the service you want to create (e.g., simple-service) 
1. Run GitHub workflow `One-time repository initialization` to tailor deployment files and go modules to the new instance of the keptn service template. This will create a Pull Request containing the necessary changes, review it, adjust if necessary and merge it.
1. Figure out whether your Kubernetes Deployment requires [any RBAC rules or a different service-account](https://github.com/keptn-sandbox/contributing#rbac-guidelines), and adapt [chart/templates/serviceaccount.yaml](chart/templates/serviceaccount.yaml) accordingly for the roles.
1. Last but not least: Remove this intro within the README file and make sure the README file properly states what this repository is about

---
# Not supported in the query
- `fillmissing`
- `outlier`
- `timeshift`

Why? Because the API does not support `fillmissing` and `outlier`. `timeshift` is supported but you can't write it in the query like `<my-query> | timeshift`. We plan to support `timeshift` in the future ([issue](https://github.com/vadasambar/sumologic-service/issues/1)) but support for `fillmissing` and `outlier` depends on Sumo Logic (can't do anything until Sumo Logic supports it). 

# Rules for using `quantize`
Based on https://help.sumologic.com/Metrics/Metric-Queries-and-Alerts/07Metrics_Operators/quantize#quantize-syntax
1. Use only 1 `quantize` (using `quantize` multiple times in a query leads to error)
2. Use `quantize` immediately after the metric query before any other operator
3. Quantize should be strictly defined as `query | quantize to [TIME INTERVAL] using [ROLLUP]` (this differs from how Sumo Logic quantize works. You need to be explicit here. Dropping [TIME INTERVAL] or `using` or `[ROLLUP]` won't work)  

Why so many rules? Because [Sumo Logic API does not support quantize in the query](https://api.sumologic.com/docs/#operation/runMetricsQueries). We have implemented a wrapper
which mimics quantize which works well if you adhere to 1-3.

# sumologic-service
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/sumologic-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/sumologic-service)](https://goreportcard.com/report/github.com/keptn-sandbox/sumologic-service)

This implements a sumologic-service for Keptn. If you want to learn more about Keptn visit us on [keptn.sh](https://keptn.sh)

## Compatibility Matrix

*Please fill in your versions accordingly*

| Keptn Version    | [sumologic-service Docker Image](https://hub.docker.com/r/keptn-sandbox/sumologic-service/tags) |
|:----------------:|:----------------------------------------:|
|       0.6.1      | keptn-sandbox/sumologic-service:0.1.0 |
|       0.7.1      | keptn-sandbox/sumologic-service:0.1.1 |
|       0.7.2      | keptn-sandbox/sumologic-service:0.1.2 |

## Installation

The *sumologic-service* can be installed as a part of [Keptn's uniform](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the *sumologic-service* in your Keptn Kubernetes cluster use the [`helm chart`](chart/Chart.yaml) file,
for example:

```console
helm install -n keptn sumologic-service chart/
```

This should install the `sumologic-service` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment sumologic-service -o wide
kubectl -n keptn get pods -l run=sumologic-service
```

### Up- or Downgrading

Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):

```console
helm upgrade -n keptn --set image.tag=$VERSION sumologic-service chart/
```

### Uninstall

To delete a deployed *sumologic-service*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:

```console
helm uninstall -n keptn sumologic-service
```

## Development

Development can be conducted using any GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode with Go plugins).

It is recommended to make use of branches as follows:

* `main`/`master` contains the latest potentially unstable version
* `release-*` contains a stable version of the service (e.g., `release-0.1.0` contains version 0.1.0)
* create a new branch for any changes that you are working on, e.g., `feature/my-cool-stuff` or `bug/overflow`
* once ready, create a pull request from that branch back to the `main`/`master` branch

When writing code, it is recommended to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).

### Where to start

If you don't care about the details, your first entrypoint is [eventhandlers.go](eventhandlers.go). Within this file 
 you can add implementation for pre-defined Keptn Cloud events.
 
To better understand all variants of Keptn CloudEvents, please look at the [Keptn Spec](https://github.com/keptn/spec).
 
If you want to get more insights into processing those CloudEvents or even defining your own CloudEvents in code, please 
 look into [main.go](main.go) (specifically `processKeptnCloudEvent`), [chart/values.yaml](chart/values.yaml),
 consult the [Keptn docs](https://keptn.sh/docs/) as well as existing [Keptn Core](https://github.com/keptn/keptn) and
 [Keptn Contrib](https://github.com/keptn-contrib/) services.

### Common tasks

* Build the binary: `go build -ldflags '-linkmode=external' -v -o sumologic-service`
* Run tests: `go test -race -v ./...`
* Build the docker image: `docker build . -t keptn-sandbox/sumologic-service:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Run the docker image locally: `docker run --rm -it -p 8080:8080 keptn-sandbox/sumologic-service:dev`
* Push the docker image to DockerHub: `docker push keptn-sandbox/sumologic-service:dev` (Note: Ensure that you use the correct DockerHub account/organization)
* Deploy the service using `kubectl`: `kubectl apply -f deploy/`
* Delete/undeploy the service using `kubectl`: `kubectl delete -f deploy/`
* Watch the deployment using `kubectl`: `kubectl -n keptn get deployment sumologic-service -o wide`
* Get logs using `kubectl`: `kubectl -n keptn logs deployment/sumologic-service -f`
* Watch the deployed pods using `kubectl`: `kubectl -n keptn get pods -l run=sumologic-service`
* Deploy the service using [Skaffold](https://skaffold.dev/): `skaffold run --default-repo=your-docker-registry --tail` (Note: Replace `your-docker-registry` with your container image registry (defaults to ghcr.io/keptn-sandbox/sumologic-service); also make sure to adapt the image name in [skaffold.yaml](skaffold.yaml))


### Testing Cloud Events

We have dummy cloud-events in the form of [RFC 2616](https://ietf.org/rfc/rfc2616.txt) requests in the [test-events/](test-events/) directory. These can be easily executed using third party plugins such as the [Huachao Mao REST Client in VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

## Automation

### GitHub Actions: Automated Pull Request Review

This repo uses [reviewdog](https://github.com/reviewdog/reviewdog) for automated reviews of Pull Requests. 

You can find the details in [.github/workflows/reviewdog.yml](.github/workflows/reviewdog.yml).

### GitHub Actions: Unit Tests

This repo has automated unit tests for pull requests. 

You can find the details in [.github/workflows/CI.yml](.github/workflows/CI.yml).

### GH Actions/Workflow: Build Docker Images

This repo uses GH Actions and Workflows to test the code and automatically build docker images.

Docker Images are automatically pushed based on the configuration done in [.ci_env](.ci_env) and the two [GitHub Secrets](https://github.com/keptn-sandbox/sumologic-service/settings/secrets/actions)
* `REGISTRY_USER` - your DockerHub username
* `REGISTRY_PASSWORD` - a DockerHub [access token](https://hub.docker.com/settings/security) (alternatively, your DockerHub password)

## How to release a new version of this service

It is assumed that the current development takes place in the `main`/`master` branch (either via Pull Requests or directly).

Once you're ready, go to the Actions tab on GitHub, select Pre-Release or Release, and run the action.


## License

Please find more information in the [LICENSE](LICENSE) file.
