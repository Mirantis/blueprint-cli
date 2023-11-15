---
title: "Bootstrap a Kind cluster"
draft: false
weight: 2
---

1. Install `Kind`: https://kind.sigs.k8s.io/docs/user/quick-start/
2. Generate a sample blueprint file:
   ```shell
   bctl init --kind > blueprint.yaml
   ```~
   This will create a blueprints file `blueprint.yaml` with a kind cluster definition, a core ingress component and an addon. See a [sample here](#sample-blueprint-for-kind-cluster)
3. Deploy the blueprint
   ```shell
   bctl apply --config blueprint.yaml
   ```
4. Connect to the cluster:~
   ```shell
   export KUBECONFIG=./kubeconfig
   kubectl get pods
   ```
   Note: `bctl` will create a `kubeconfig` file in the current directory.
   Use this file to connect to the cluster.
5. Add wordpress addon to the `blueprint.yaml`:
   ```YAML
   - name: wordpress
     kind: HelmAddon
     enabled: true
     namespace: wordpress
     chart:
       name: wordpress
       repo: https://charts.bitnami.com/bitnami
       version: 18.0.11
   ```
   Update your cluster with the updated blueprint:

   ```shell
   bctl update --config blueprint.yaml
   ```
6. Delete the cluster:
   ```shell
   bctl reset --config blueprint.yaml
