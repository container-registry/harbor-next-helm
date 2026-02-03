# Deploying Harbor on Openshift

## Set URL

externalURL must be replaced by the public URL provided by Openshift. In our case it is:

```yaml
externalURL: "https://alex-container-regis-dev.apps.rm2.thpm.p1.openshiftapps.com"
```

## Ingress

Openshift needs an annotation for routing to work. We add the termination type:

```yaml
  annotations:
    # cert-manager.io/cluster-issuer: letsencrypt-prod
    route.openshift.io/termination: edge
```

Hosts and core must be set to the URL provided by Openshift. For example:

```yaml
  core: "alex-container-regis-dev.apps.rm2.thpm.p1.openshiftapps.com"
  # -- Additional ingress hosts
  hosts:
    - host: alex-container-regis-dev.apps.rm2.thpm.p1.openshiftapps.com
      paths:
        - path: /
          pathType: Prefix
```

## Security contexts

`securityContext` and `podSecurityContext` are automatically set during deployment so we initialize it to empty objects:

```yaml
  securityContext: {}
  podSecurityContext: {}
```

In `valkey` section, we set `securityContext` and `podSecurityContext` to null values:

```yaml
  securityContext:
    readOnlyRootFilesystem: null
    runAsNonRoot: null
    runAsUser: null

  podSecurityContext:
    fsGroup: null
    runAsUser: null
    runAsGroup: null
```

## Storage class

`storageClass` fields are set to `gp3` in each persistence section (default on Openshift):

```yaml
persistence:
  enabled: true
  storageClass: "gp3"
```

# Deploying a postgresql database for Harbor

To have a working Harbor deployment, we can install any Postgresql database and
reference it in the values.

# Access Harbor

Harbor, after a successful deployment, should be available on the `externalURL`
provided. If it does not work, check the events on the dashboard of Openshift.
