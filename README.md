# event-cwl-exporter

## ビルド

macOS のみ対応しています。

``` shellsession
make
```

## イメージの Push

``` shellsession
aws ecr create-repository --repository-name event-cwl-exporter
docker build . -t [URI]:[Version]
aws ecr get-login-password | docker login --username AWS --password-stdin [ECR endpoint]
docker push [URI]:[Version]
```

## デプロイ

デフォルトリージョンは us-west-2 になっています。

***yaml/event-cwl-exporter.yaml*** の以下のプレースホルダを置換してください。

* {{EKS cluster name}}
* {{Image URI}}

その上で以下のコマンドを実行します。

``` shellsession
kubectl apply -f yaml/event-cwl-exporter.yaml
```

なお、ワーカーノードのロールに CloudWatchAgentServerPolicy のポリシーが必要です。
