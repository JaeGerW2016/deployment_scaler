### Build image

```
docker build -t 314315960/deployment_scaler .

```

### Create CronJob On Kubernetes
```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: scale-myapp
  namespace: kube-system
spec:
  schedule: "23 2 * * *"
  jobTemplate:
    metadata:
      name: scale-myapp
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: "true"
    spec:
      template:
        spec:
          serviceAccountName: admin-user
          containers:
          - name: scale-myapp
            image: 314315960/deployment_scaler:latest
            command:
            - /root/deployment_scaler
            - --deployment
            - myapp
            - --replicas
            - '5'
          restartPolicy: "OnFailure"

```
