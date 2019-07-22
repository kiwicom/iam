# End to End tests

These tests create the Kiwi IAM service and required services (Okta nginx, Redis) and runs tests against it to guarantee that our service is running correctly.

## Running locally

To run e2e tests you will need to have `docker-compose` installed.

```sh
make e2e
```

If you didn't make changes on IAM or mocks (eg. you only created e2e tests) you
can run `make e2e/nobuild` which will skip the image builds.

Both `make e2e` and `make e2e/nobuild` start a test environment using
docker-compose, to tear down this environment run `e2e/env-stop`.
