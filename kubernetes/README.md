This directory has the needed resource to deploy quizmaker on Kubernetes

TODO:
- [ ] Create the needed resources (deployment, service, ingress etc)
- [ ] Create a kustomization to be used remotely
- [ ] Create a kairos config that can be used to automatically deploy quizmaker
- [ ] Write instructions here on how to use the above resources

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

secretGenerator:
  - name: quizmaker-secret
    literals:
      - token=<your_cookie_secret_here>

configMapGenerator:
  - name: quizmaker-questions
    files:
      - questions.yaml
```

- Apply that to your cluster with:

```
kubectl apply -k overlays
```
