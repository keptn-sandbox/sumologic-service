# sumologic-service
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/sumologic-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/sumologic-service)](https://goreportcard.com/report/github.com/keptn-sandbox/sumologic-service)

This implements the `sumologic-service` that integrates the [Sumo Logic](https://en.wikipedia.org/wiki/Sumo_Logic) observability platform with Keptn. This enables you to use Sumo Logic as the source for the Service Level Indicators ([SLIs](https://keptn.sh/docs/0.15.x/reference/files/sli/)) that are used for Keptn [Quality Gates](https://keptn.sh/docs/concepts/quality_gates/).
If you want to learn more about Keptn visit us on [keptn.sh](https://keptn.sh)

You can find more information about the service in the [proposal issue](https://github.com/keptn/integrations/issues/20)

## Quickstart
If you are on Mac or Linux, you can use [examples/kup.sh](./examples/kup.sh) to set up a local Keptn installation that uses Sumo Logic. This script creates a local minikube cluster, installs Keptn, Istio, Sumo Logic and the Sumo Logic integration for Keptn (check the script for pre-requisites). 

To use the script,
```bash
export ACCESS_ID="<your-sumologic-access-id>" ACCESS_KEY="<your-sumologic-access-key>"
examples/kup.sh
```
Check [the official docs](https://help.sumologic.com/Manage/Security/Access-Keys#manage-your-access-keys-on-preferences-page) for how to create the Sumo Logic access ID and access key

## If you already have a Keptn cluster running
**1. Install Sumo Logic**
```bash
export ACCESS_ID="<your-sumologic-access-id>" ACCESS_KEY="<your-sumologic-access-key>"
helm upgrade --install my-sumo sumologic/sumologic   --set sumologic.accessId="${ACCESS_ID}"   --set sumologic.accessKey="${ACCESS_KEY}"   --set sumologic.clusterName="keptn-sumo"

```
**2. Install Keptn sumologic-service to integrate Sumo Logic with Keptn**
```bash
export ACCESS_ID="<your-sumologic-access-id>" ACCESS_KEY="<your-sumologic-access-key>"
# cd sumologic-service
helm install sumologic-service ../helm --set sumologicservice.accessId=${ACCESS_ID} --set sumologicservice.accessKey=${ACCESS_KEY} 

```

**3. Add SLI and SLO**
```bash
keptn add-resource --project="<your-project>" --stage="<stage-name>" --service="<service-name>" --resource=/path-to/your/sli-file.yaml --resourceUri=sumologic/sli.yaml
keptn add-resource --project="<your-project>"  --stage="<stage-name>" --service="<service-name>" --resource=/path-to/your/slo-file.yaml --resourceUri=slo.yaml
```
Example:
```bash
keptn add-resource --project="podtatohead" --stage="hardening" --service="helloservice" --resource=./quickstart/sli.yaml --resourceUri=sumologic/sli.yaml
keptn add-resource --project="podtatohead" --stage="hardening" --service="helloservice" --resource=./quickstart/slo.yaml --resourceUri=slo.yaml
```
Check [./quickstart/sli.yaml](./examples/quickstart/sli.yaml) and [./quickstart/slo.yaml](./examples/quickstart/slo.yaml) for example SLI and SLO. 

<!-- TODO: Uncomment this after the PR to support switching SLI provider is merged -->
<!-- 4. Configure Keptn to use Sumo Logic SLI provider
Use keptn CLI version [0.15.0](https://github.com/keptn/keptn/releases/tag/0.15.0) or later.
```bash
keptn configure monitoring sumologic --project <project-name>  --service <service-name>
``` -->
**4. Configure Keptn to use Sumo Logic SLI provider**  

There's an [open PR](https://github.com/keptn/keptn/pull/8546) to support `keptn configure monitoring sumologic` in the future releases but for now, you need to configure Keptn to use Sumo Logic manually by creating a ConfigMap like this:
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lighthouse-config-<your-project-name>
  namespace: keptn
data:
  sli-provider: "sumologic"

```
[Example](./examples/quickstart/lighthouse_config.yaml)
```
kubectl apply -f <above-configmap-file>
```

**5. Trigger delivery**
```bash
keptn trigger delivery --project=<project-name> --service=<service-name> --image=<image> --tag=<tag>
```
Example:
```bash
keptn trigger delivery --project=podtatohead --service=helloservice --image=docker.io/jetzlstorfer/helloserver --tag=0.1.1
```
Observe the results in the [Keptn Bridge](https://keptn.sh/docs/0.15.x/bridge/)

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
which mimics quantize which works well if you adhere to the above rules.

## Compatibility Matrix

| Keptn Version    | [sumologic-service Docker Image](https://github.com/keptn-sandbox/sumologic-service/pkgs/container/sumologic-service) | Sumo Logic chart |
|:----------------:|:----------------------------------------:| :----------------------------------------: |
|       0.15.0      | keptn-sandbox/sumologic-service:0.15.0 | [2.14.1](https://github.com/SumoLogic/sumologic-kubernetes-collection/tree/v2.14.1/deploy/helm/sumologic) | 

## Installation

```bash
export ACCESS_ID="<your-sumologic-access-id>" ACCESS_KEY="<your-sumologic-access-key>"
# cd sumologic-service
helm upgrade --install my-sumo sumologic/sumologic   --set sumologic.accessId="${ACCESS_ID}"   --set sumologic.accessKey="${ACCESS_KEY}"   --set sumologic.clusterName="keptn-sumo"

```
<!-- TODO: Uncomment this after the PR to support switching SLI provider is merged -->
<!-- Tell Keptn to use Sumo Logic as SLI provider for your project/service
```bash
keptn configure monitoring sumologic --project <project-name>  --service <service-name>
``` -->

Tell Keptn to use Sumo Logic as the SLI provider for your project/service ([future releases will support a better way to do this](https://github.com/keptn/keptn/pull/8546)):
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lighthouse-config-<your-project-name>
  namespace: keptn
data:
  sli-provider: "sumologic"

```
[Example](./examples/quickstart/lighthouse_config.yaml)
```
kubectl apply -f <above-configmap-file>
```

This should install the `sumologic-service` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment sumologic-service -o wide
kubectl -n keptn get pods -l run=sumologic-service
```

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

```bash
helm upgrade sumologic-service ./helm --set sumologicservice.accessId=${ACCESS_ID} --set sumologicservice.accessKey=${ACCESS_KEY} 
```

### Uninstall

To delete a deployed *sumologic-service* helm chart:

```bash
helm uninstall sumologic-service
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
* Build the docker image: `make docker-build IMG_TAG=dev`
* Run the docker image locally: `docker run --rm -it -p 8080:8080 ghcr.io/keptn-sandbox/sumologic-service:dev`
* Push the docker image to DockerHub: `docker push ghcr.io/keptn-sandbox/sumologic-service:latest` 
* Watch the deployment using `kubectl`: `kubectl -n keptn get deployment sumologic-service -o wide`
* Get logs using `kubectl`: `kubectl -n keptn logs deployment/sumologic-service -f`
* Watch the deployed pods using `kubectl`: `kubectl -n keptn get pods -l run=sumologic-service`


### Testing Cloud Events

We have dummy cloud-events in the form of [RFC 2616](https://ietf.org/rfc/rfc2616.txt) requests in the [test-events/](test-events/) directory. These can be easily executed using third party plugins such as the [Huachao Mao REST Client in VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

## Automation

### GitHub Actions: Automated Pull Request Review

This repo uses [reviewdog](https://github.com/reviewdog/reviewdog) for automated reviews of Pull Requests. 

You can find the details in [.github/workflows/reviewdog.yml](.github/workflows/reviewdog.yml).

### GitHub Actions: Unit Tests

This repo has automated unit tests for pull requests. 

You can find the details in [.github/workflows/CI.yml](.github/workflows/CI.yml).

## How to release a new version of this service

It is assumed that the current development takes place in the `main`/`master` branch (either via Pull Requests or directly).

Once you're ready, go to the Actions tab on GitHub, select Pre-Release or Release, and run the action.


## License

Please find more information in the [LICENSE](LICENSE) file.
