# 280blocker_proxy
## 何がうれしい？
年月の部分を自動補完します。AdGuard Homeなどに登録するときに、手動更新が不要になるメリットがあります。

## Usage
### access
```
curl 127.0.0.1:3030/healthz
{"status":"ok","time":"2025-11-08T13:47:37Z"}
```
```
curl 127.0.0.1:3030/280blocker.txt
...
```
### test
```
go test ./...
```

### run (docker)
```
docker compose up --build -d
```

### run (k8s)
```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bp280-server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: bp280-server
  template:
    metadata:
      labels:
        app: bp280-server
    spec:
      containers:
        - name: bp280-server
          image: alicey/280blocker_proxy:latest
          ports:
            - containerPort: 36745
          env:
            - name: PORT
              value: "36745"
          securityContext:
            readOnlyRootFilesystem: true
          readinessProbe:
            httpGet:
              path: /healthz
              port: 36745
            initialDelaySeconds: 30
            timeoutSeconds: 10
            periodSeconds: 60
            failureThreshold: 3
          livenessProbe:
            httpGet:
              path: /healthz
              port: 36745
            initialDelaySeconds: 60
            timeoutSeconds: 10
            periodSeconds: 60
            failureThreshold: 5
```