# Kiwi IAM

## Usage

Create `.env.yaml` file and set environment variables. Check `.env-sample.yaml`
to see all the possible variables.

To install dependencies and start the project, make sure your GOPATH is set,
and run:

```
make
make start
```

You can use `make dev` on development to reload the server automatically on file
changes.

# Redis

To run this project you need to have redis installed and running, you can use
the following commands to do so.

```shell
# MacOS
brew install redis
brew services start redis

# Linux
# install using your package manager
systemctl start redis

# Docker
docker run -it --rm -p 6379:6379 redis
```

You can use `redis-cli` to interact with redis (ie. check keys and values).
Useful commands:

```shell
# Show all keys
KEYS *

# Show value for a key
GET <key>
```
# API Documentation

We use [`go-swagger`](https://github.com/go-swagger/go-swagger) for generating the API documentation based on the OpenAPI 2.0 specification which is defined on `api/swagger.yml`.

```shell
# Opens the documentation website locally (on http://localhost:51261/docs)
make swagger-serve

# Validates the swagger schema
make swagger-validate
```

# Secret Manager

This service uses Vault for syncing secrets to our app.

The following structure is expected:

- `secret/{VAULT_NAMESPACE}/{environment}/app_tokens`
- `secret/{VAULT_NAMESPACE}/{environment}/settings`

## Contributing

- Run `make test/all` before pushing changes.
- Commit messages should be at most 72 characters long.
- Commit messages should start with the scope of the changes introduced by the commit.

## Contributing to documentation

Run `make install_hugo` to install hugo, and `make docs/serve` to launch the
hugo server for development.

## Troubleshooting

### It cannot find go or some packages when executing some make commands

When executing `swagger`, or `protoc`, or running `make dev` (after having run `make`), it says that those command do not exist...

In golang, the main packages are installed in your go installation `bin` folder, which defaults to `usr/local/go/bin`. When installing new packages through `make`, the binaries are installed in `$GOPATH/bin`. ([How to set GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)).

To fix the issue you have make them accessible through the terminal, by adding both of these routes to your `$PATH` so they are known to your bash environment, i.e editing your `.bashrc` or `.zshrc`

```shell
export GOPATH=$HOME/go #example of GOPATH directory
export PATH=$PATH:/usr/local/go/bin $PATH:$GOPATH/bin
``

