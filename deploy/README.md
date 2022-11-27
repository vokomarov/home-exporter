# Infrastructure

Root directory of infrastructure related files, service deployment and configuration.

## Production

Namespace: `telegram-bots`

### Configure

```shell
$ cp config.yml.example config.yml
$ cp secret.yml.example secret.yml
```

Set configuration keys and tokens in `secret.yml` and `config.yml` files

### Install

```shell
$ kubectl apply -f namespace.yml -f config.yml -f secret.yml -f service.yml -f deployment.yml
```

### Deploy

```shell
$ kubectl -n telegram-bots set image deployment/home-exporter home-exporter=vovanms/home_exporter:0.0.1      # Deploy new tag
$ kubectl -n telegram-bots rollout status deployment/home-exporter                                           # Watch deployment status
$ kubectl -n telegram-bots rollout undo deployment/home-exporter                                             # Rollback current deployment
$ kubectl -n telegram-bots rollout history deployment/home-exporter                                          # List past deployment revision
$ kubectl -n telegram-bots rollout restart deployment/home-exporter                                          # Redeploy currently deployed tag
```


