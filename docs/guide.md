# Create a Harbor deployment on rancher

## Supported versions

Rancher manager >= 2.13.1

RKE2 >= rke2 version v1.34.3+rke2r1 (1b103f296ab20fac6b32951c9efe59d28a5ed79f)

## Deploy a database

We create a namespace where our database will be running:

```bash
kubectl create namespace my-db
```

In this guide we use local storage, so we install a storage provisioner:

```bash
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.34/deploy/local-path-storage.yaml
```

Set it as default:

```bash
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```

We then create this storage class and set it as default:

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: local
provisioner: rancher.io/local-path
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

We deploy our external Postgresql database with the Bitnami Helm chart:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgresql bitnami/postgresql \
  --namespace my-db \
  --set auth.postgresPassword=test1234! \
  --set persistence.size=8Gi
```

Once installed, let us connect to it by forwarding the port 5432 to the host:

```bash
kubectl port-forward --namespace=my-db postgresql-0 5432:5432


Forwarding from 127.0.0.1:5432 -> 5432
Forwarding from [::1]:5432 -> 5432
```

Then we connect to our database from the host by entering the password given
during installation:

```bash
psql -h localhost -U postgres -p 5432
Password for user postgres:
```

We create a `registry` database which will be used by Harbor:

```sql
postgres=# create database registry;
CREATE DATABASE
postgres=# exit
```

## Deploy Harbor next

We pull the Helm chart and decompress it:

```bash
helm pull https://8gears.container-registry.com/harbor-next/chart/harbor
tar xzvf harbor-x.x.x.tgz
```

We update `values.yaml` to set the hostname, port and credentials of our
database we will use with Harbor:

```yaml
database:
  host: "postgresql.my-db.svc.cluster.local"
  port: 5432
  username: "postgres"
  password: "test1234!"
  database: registry
```

The hostname chosen for our local deployment is the default,
`harbor.example.com`. To override our DNS resolver, we set a mapping to localhost in
`/etc/hosts`:

```
127.0.0.1 harbor.example.com
```

We will now deploy Harbor:

```bash
helm upgrade --install test-1 . \
  --namespace my-container-registry \
  --create-namespace
```

After some time we now see all the pods and our installation running:

```bash
kubectl get pods -n my-container-registry
NAME                                        READY   STATUS    RESTARTS      AGE
test-1-harbor-core-754c57dd77-xqfkb         1/1     Running   0             61s
test-1-harbor-exporter-867746db74-5pgmn     1/1     Running   0             61s
test-1-harbor-jobservice-67bc99bd86-7bvlp   0/1     Running   3 (35s ago)   61s
test-1-harbor-portal-57f7894ff7-hfs8h       1/1     Running   0             60s
test-1-harbor-registry-5bd8d59648-9fqs5     1/1     Running   0             61s
test-1-harbor-scanner-trivy-0               1/1     Running   0             61s
test-1-valkey-59486f6977-nqb2z              1/1     Running   0             61s
```

## Access Harbor from an Nginx ingress

We will use an Nginx ingress to forward traffic to our host computer. Let us install it with Helm:

```bash
helm upgrade --install ingress-nginx ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace ingress-nginx --create-namespace
```

Given the ingress is created:

```bash
kubectl get ingress -A
NAMESPACE               NAME            CLASS    HOSTS                ADDRESS   PORTS     AGE
my-container-registry   test-1-harbor   <none>   harbor.example.com             80, 443   21m
```

And endpoints are created as well:

```bash
kubectl get endpoints -n my-container-registry
Warning: v1 Endpoints is deprecated in v1.33+; use discovery.k8s.io/v1 EndpointSlice
NAME                          ENDPOINTS                         AGE
test-1-harbor-core            10.42.0.38:8080,10.42.0.38:8001   22m
test-1-harbor-exporter        10.42.0.39:8001                   22m
test-1-harbor-jobservice      10.42.0.40:8080,10.42.0.40:8001   22m
test-1-harbor-portal          10.42.0.41:8080                   22m
test-1-harbor-registry        <none>                            22m
test-1-harbor-scanner-trivy   <none>                            22m
test-1-valkey                 10.42.0.34:6379                   22m
```

We should now be able to access our deployment by forwarding the ingress port:

```
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 4443:443
Forwarding from 127.0.0.1:4443 -> 443
Forwarding from [::1]:4443 -> 443
Handling connection for 4443
```


We can access the portal from the browser at

```
https://harbor.example.com:4443
```

The default credentials are `admin` and `Harbor12345`.
