# event-cwl-exporter


```
aws ecr create-repository --repository-name event-cwl-exporter
docker build . -t [URI]:[Version]
aws ecr get-login-password | docker login --username AWS --password-stdin [ECR endpoint]
```
