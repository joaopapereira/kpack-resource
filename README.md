# kpack Concourse Resource

This Concourse Resource allow the user to schedule builds directly from
the pipeline

## Source Configuration

* `a`: *Required.* This is a required setting.

* `b`: *Optional.* This is an optional setting.

### Example

```yaml
resource_types:
- name: kpack
  type: docker-image
  source:
    repository: joaopapereira/kpack-resource

resources:
- name: first_image
  type: kpack
  check_every: 5m
  source:
    k8s_token: "token from a service account in k8s"
    k8s_api: "https://k8s.api:port"
    k8s_ca_cert: "CA Cert for k8s api"
    image: "gcr.io/some/image:tag"
    username: "gcr username to Get the image"
    password: "gcr password"
    namespace: "namespace-in-k8s-where-images-will-be-created"
    service_account: "service-account-that-as-credentials-to-push-image"
    log_level: debug
- name: kpack
  type: git
  source:
    uri: https://github.com/pivotal/kpack
    branch: master


jobs:
- name: Produce Image
  plan:
  - get: kpack
  - put: first_image
    params:
      git_path: kpack
```

## Behavior

### `check`: Check for something

Checks if there is a new version of the image in K8s.

### `in`: Fetch something

Same as the docker resource. Pulls the image created by kpack

#### Parameters

### `out`: Put something somewhere

Creates or updates an image in kubernetes.

#### Parameters

* `git_path`: *Optional. Path in which the git Resource downloaded the source to

* `blob_path`: *Optional. URL to the blob that can be used to create a new build from.


**Note:** At least one of the above need to be set

## Development

### Prerequisites

* golang is *required* - version 1.13.x or higher is required.
* docker is *required* - version 17.05.x or higher is required.
* make is *required* - version 4.1 of GNU make is tested.

### Running the tests

The Makefile includes a `test` target, and tests are also run inside the Docker build.

Run the tests with the following command:

```sh
make test
```

### Building and publishing the image

The Makefile includes targets for building and publishing the docker image. Each of these
takes an optional `VERSION` argument, which will tag and/or push the docker image with
the given version.

```sh
make VERSION=1.2.3
make publish VERSION=1.2.3
```
