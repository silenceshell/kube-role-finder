## kube-role-finder

Get clusterRole which has define the specified resource.

## build

Just run the `build.sh`.

## donwload

- [Mac版](https://silenceshell-1255345740.cos.ap-shanghai.myqcloud.com/kube-role-finder)
- [Linux版](https://silenceshell-1255345740.cos.ap-shanghai.myqcloud.com/kube-role-finder-linux)

## usage

### find all clusterRole which has defined resource `services`

```
kube-role-finder -resource services
```

### find all clusterRole which has defined resource `apps/deployments`

```
kube-role-finder -apiGroup apps -resource deployments
```

### find all clusterRole which has defined resource `apps/deployments` with verb "delete"

```
kube-role-finder -apiGroup apps -resource deployments -verb 
```
