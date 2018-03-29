
Setup
====

1. Run the server
```
dep ensure
go run main.go --assistant <IP of Google Home> --port 8080
```

2. Export local http port
```
ngrok http 8080
```

3. Set up LINE Message's webhook
Set webhook URL [via web console.](https://developers.line.me/console/channel/1570827969/basic/)



