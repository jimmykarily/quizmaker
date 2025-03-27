This directory has the needed resource to deploy quizmaker on Kubernetes

## Deploy

- Create a directory for your overlays:

```
mkdir overlays
```

- Create your own questions.yaml file

```
cat << EOF > overlays/questions.yaml
<put_your_questions_yaml_here>
EOF
```

- Create a kustomization file:

```
cat << EOF > overlays/kustomization.yaml
resources:
  - https://github.com/jimmykarily/quizmaker/kubernetes?ref=main

patchesStrategicMerge:
  - domain-patch.yaml
  - cluster-issuer-patch.yaml

secretGenerator:
  - name: quizmaker-secret
    behavior: replace
    literals:
      - token="<your_secret_here>"
      - verificationEndpoint="<a_random_url_path_here>"

configMapGenerator:
  - name: quizmaker-questions
    behavior: replace
    files:
      - questions.yaml

generatorOptions:
  disableNameSuffixHash: true
EOF
```

- Create the 2 patch files to set the cluster issuer email and your domain
  for the ingress:

```
cat << EOF > overlays/cluster-issuer-patch.yaml
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
  namespace: cert-manager
spec:
  acme:
    email: <your_email_here>

---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-production
  namespace: cert-manager
spec:
  acme:
    email: <your_email_here>
EOF
```

```
cat << EOF > overlays/domain-patch.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: quizmaker-ingress
spec:
  rules:
    - host: <your_domain_here>
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: quizmaker-service
                port:
                  number: 80
  tls:
    - hosts:
        - <your_domain_here>
      secretName: quizmaker-tls
```

- Apply that to your cluster with:

```
kubectl apply -k overlays/
```


## NOTES

If you are deploying on Kairos, you have to make sure k3s is deployed with additional
SANs matching your domain. Here is an example Kairos config that can be used
as userdata in AWS:

```
#cloud-config

users:
  - name: kairos
    ssh_authorized_keys:
      - github:<your_handle_here>
    groups:
      - admin

k3s:
  enabled: true
  args:
  - --tls-san "<your_domain_here>"

reset:
  system:
    uri: "quay.io/kairos/opensuse:leap-15.6-standard-amd64-generic-v3.4.0-beta7-k3s1.31.6-k3s1"
```
