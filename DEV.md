
Running locally
go run .

Deployment
O111MODULE=on gcloud --project=go-tour-cz app deploy --no-promote app.yaml

Removing previous app versions
```
gcloud --project=go-tour-cz app versions list
gcloud --project=go-tour-cz app versions delete [VERSION_ID]
```

See GCloud logs 
```
gcloud app logs tail -s default
```


