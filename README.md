# Compose specification compliance test suite

This repository contains a test suite for testing Compose implementations to
ensure that they correctly implement the Compose specification. The current
state of the test suite is that it is a work in progress and contributions are
welcome!

## Getting started

By default, the test suite is run against the
[Compose reference implementation](https://github.com/compose-spec/compose-ref)
which uses the Docker Engine.

### Prerequisites

* [Docker](https://docs.docker.com/install/)
* [compose-ref](https://github.com/compose-spec/compose-ref) in your PATH
* Ensure that you have no running containers (see [issue](https://github.com/compose-spec/compatibility-test-suite/issues/5))

### Running the tests

Using the defaults, you can run the tests as follows:

```console
$ make test
```
