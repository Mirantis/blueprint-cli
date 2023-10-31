# Boundless Operator - Tech Preview

<!-- TOC -->
* [Quick Start](#quick-start)
  * [Install on Kind](#install-on-kind)
  * [Install on Amazon VM](#install-on-amazon-vm)
* [Boundless Operator Blueprints](#boundless-operator-blueprints)
  * [Core Components](#core-components)
  * [Add-ons](#add-ons)
* [Sample Blueprints](#sample-blueprints)
<!-- TOC -->

## Quick Start

### Install on Kind

1. Install `Kind`: https://kind.sigs.k8s.io/docs/user/quick-start/
2. Install Boundless CLI binary:
   ```shell
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mirantis/boundless/main/script/install.sh)"
   ```
   This will install `bctl` to `/usr/local/bin`. See [here](https://github.com/Mirantis/boundless/releases) for all releases.
3. Generate a basic blueprint file:
   ```shell
   bctl init --kind > blueprint.yaml
   ```
   This will create a basic blueprints file `blueprint.yaml`. See a [sample here](#sample-blueprint-for-kind-cluster)
4. Create the cluster:
   ```shell
   bctl apply --config blueprint.yaml
   ```
5. Connect to the cluster:
   ```shell
   export KUBECONFIG=./kubeconfig
   kubectl get pods
   ```
   Note: `bctl` will create a `kubeconfig` file in the current directory.
   Use this file to connect to the cluster.
6. Update the cluster by modifying `blueprint.yaml` and then running:
   ```shell
   bctl update --config blueprint.yaml
   ```
7. Delete the cluster:
   ```shell
   bctl reset --config blueprint.yaml
   ```

### Install on Amazon VM

### Prerequisites
Ensure that following are installed on the system:
* `k0sctl` (required for installing k0s distribution): https://github.com/k0sproject/k0sctl#installation
* `terraform` (for creating VMs on AWS)

### Create virtual machines on AWS

There are `terraform` scripts in the `example/` directory that can be used to create machines on AWS.

1. `cd example/aws-tf`
2. Create a `terraform.tfvars` file with the content similar to:
   ```
   cluster_name = "rs-boundless-test"
   controller_count = 1
   worker_count = 1
   cluster_flavor = "m5.large"
   ```
3. `terraform init`
4. `terraform apply`
5. `terraform output --raw bop_cluster > ./blueprint.yaml`

### Install Boundless Operator on `k0s`

1. Install Boundless CLI binary:
   ```shell
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mirantis/boundless/main/script/install.sh)"
   ```
   This will install `bctl` to `/usr/local/bin`. See [here](https://github.com/Mirantis/boundless/releases) for all releases.
2. Generate a basic blueprint file:
   ```shell
   bctl init > blueprint.yaml
   ```
   This will create a basic blueprints file `blueprint.yaml`. See a [sample here](#sample-blueprint-for-k0s-cluster)
3. Now, edit the `blueprint.yaml` file to set the `spec.infra.hosts` from the output of `terraform output --raw bop_cluster`.

   The host section should look similar to:
   ```yaml
   spec:
     infra:
       hosts:
       - ssh:
           address: 52.91.89.114
           keyPath: ./example/aws-tf/aws_private.pem
           port: 22
           user: ubuntu
         role: controller
       - ssh:
           address: 10.0.0.2
           keyPath: ./example/aws-tf/aws_private.pem
           port: 22
           user: ubuntu
         role: worker
   ```
4. Create the cluster:
   ```shell
   bctl apply --config blueprint.yaml
   ```
5. Connect to the cluster:
   ```shell
   export KUBECONFIG=./kubeconfig
   kubectl get pods
   ```
   Note: `bctl` will create a `kubeconfig` file in the current directory.
   Use this file to connect to the cluster.
6. Update the cluster by modifying `blueprint.yaml` and then running:
   ```shell
   bctl update --config blueprint.yaml
   ```
7. Delete the cluster:
   ```shell
   bctl reset --config blueprint.yaml
   ```
8. Delete virtual machines:
   ```bash
   cd example/aws-tf
   terraform destroy --auto-approve
   ```

## Boundless Operator Blueprints

### Core Components

Currently, you can replace the ingress controller from `ingress-nginx` to `kong` by updating the `blueprint.yaml` file:
```yaml
spec:
 components:
   core:
     ingress:
       enabled: true
       provider: kong # ingress-nginx, kong, etc.
```

> If the cluster is already deployed, run `bctl reset` to destroy the cluster and then run `bctl apply` to recreate it.

### Add-ons
Update the `blueprint.yaml` file to add add-ons to the cluster. The add-ons are defined in the `spec
.components.addons` section.

Any public Helm chart can be used as an add-on.

Use the following configuration to add the `grafana` as an add-on:
```yaml
spec:
 components:
   addons:
   - name: my-grafana
     enabled: true
     kind: HelmAddon
     namespace: monitoring
     chart:
       name: grafana
       repo: https://grafana.github.io/helm-charts
       version: 6.58.7
       values: |
         ingress:
           enabled: true
```
and then run `bctl update` to update the cluster.

## Sample Blueprints

### Sample Blueprint for `Kind` cluster:
```yaml
apiVersion: bctl.mirantis.com/v1alpha1
kind: Cluster
metadata:
  name: kind-cluster
spec:
  kubernetes:
    provider: kind
  components:
    core:
      ingress:
        enabled: true
        provider: ingress-nginx
        config:
          controller:
            service:
              nodePorts:
                http: 30000
                https: 30001
              type: NodePort
    addons:
      - name: example-server
        kind: HelmAddon
        enabled: true
        namespace: default
        chart:
          name: nginx
          repo: https://charts.bitnami.com/bitnami
          version: 15.1.1
          values: |
            "service":
              "type": "ClusterIP"

```

### Sample Blueprint for `k0s` cluster:

```yaml
apiVersion: bctl.mirantis.com/v1alpha1
kind: Cluster
metadata:
  name: bop-cluster
spec:
  infra:
    hosts:
    - ssh:
        address: 52.91.89.114
        keyPath: ./example/aws-tf/aws_private.pem
        port: 22
        user: ubuntu
      role: controller
    - ssh:
        address: 10.0.0.2
        keyPath: ./example/aws-tf/aws_private.pem
        port: 22
        user: ubuntu
      role: worker
  kubernetes:
    provider: k0s
    version: 1.27.4+k0s.0
  components:
   core:
     ingress:
       enabled: true
       provider: ingress-nginx
       config:
         controller:
           service:
             nodePorts:
               http: 30000
               https: 30001
             type: NodePort
   addons:
   - name: example-server
     kind: HelmAddon
     enabled: true
     namespace: default
     chart:
       name: nginx
       repo: https://charts.bitnami.com/bitnami
       version: 15.1.1
       values: |2
         "service":
           "type": "ClusterIP"
```










