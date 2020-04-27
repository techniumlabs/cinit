# cinit ![cinit](https://github.com/techniumlabs/cinit/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/techniumlabs/cinit/badge.svg?branch=master)](https://coveralls.io/github/techniumlabs/cinit?branch=master)
cinit is a minimal init system for containers. It is designed to run as PID 1 inside containers. It is integrated with vault

## Why you need an init system

I can't explain it better than Yelp in *dumb-init* repo, so please [read their explanation](https://github.com/Yelp/dumb-init/blob/v1.2.0/README.md#why-you-need-an-init-system)

Summary:
- Proper signal forwarding
- Orphaned zombies reaping


## Why not just go-init
Most enterprise companies use vault or some secret store. We need some way to get those secrets into the container. That is why `cinit` exists.

## Features
1. Secret store integration.
   - Vault
   - AWS Secret Manager
2. Template substitution
   - Most of the time the resulting secrets has to be injected into some config file. We now support go templates.

## How to run
### Dockerfile
you can use it in dockerfile as follows

```
FROM alpine

USER <user>
WORKDIR /home/<user>

COPY cinit.yaml /home/<user>/.cinit.yaml

ENTRYPOINT ["cinit", "--"]
```

### Locally
    `cinit -- command_to_run`

## Configuration
cinit looks for configuration file `.cinit.yaml` in the current directory or home directory in that order. Example Sample file looks like this.
```
provider:
    secret:
        - vault

templates:
    - src: some-template-file-path
      dest: dest-path
```

To save the env variable to a file, prefix `env:` to the `src` followed by name of env variable.
```
templates:
    - src: env:MY_ENV_VAR
      dest: dest-path
```

## Acknowledgement
The code for init is taken from [go-init](https://github.com/pablo-ruth/go-init).
