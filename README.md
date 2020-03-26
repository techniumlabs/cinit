# cinit
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
   - **TODO** AWS Secret Manager
2. Template substitution
   - Most of the time the resulting secrets has to be injected into some config file. We now support go templates.

## How to run
### Dockerfile
you can use it as an entrypoint and command
    `ENTRYPOINT ["cinit", "run" "--"]`
    `CMD ["command_to_run"]`

### Locally
    `cinit run -- command_to_run`

## Configuration
cinit looks for configuration file `.cinit.yaml` in the current directory or home directory in that order.

## Acknowledgement
The code for init is taken from [go-init](https://github.com/pablo-ruth/go-init).
