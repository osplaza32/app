## Requirements

* Working Docker Desktop install with Kubernetes enabled
* [`docker-app` with CNAB support](https://github.com/docker/app/releases/tag/cnab-dockercon-preview) installed
* Source code from this directory
* Create a context with `docker-app context`
* Set the `DOCKER_TARGET_CONTEXT` environment variable

## Examples

Show the details of the application with `inspect`

```console
$ docker-app inspect
hello 0.1.0

Maintained by: garethr <garethr@docker.com>

sample app for DockerCon

Service (1) Replicas Ports Image
----------- -------- ----- -----
hello       1        8765  hashicorp/http-echo:latest

Parameters (3) Value
-------------- -----
port           8765
text           hello DockerCon
version        latest
```

Install the application:

```console
$ docker-app install
```

Show the details of the installation:

```console
$ docker-app status hello
```

Upgrade the installation, demonstrating setting parameters:

```console
$ docker-app upgrade --set port=9876 --set text="hello DockerCon EU" hello
```

Uninstall the application installation:

```console
$ docker-app uninstall hello
```

Demonstrate building a `bundle.json` for CNAB.

```console
$ docker-app bundle
Invocation image "hello:0.1.0-invoc" successfully built
```
