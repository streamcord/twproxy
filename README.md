# twproxy

Simple drop-in proxy server for the Twitch API.
Built for use with microservice architecture in mind.

## Features

- Shared ratelimit handling
- Automatic bearer token refreshing
- Multiple Client ID / Secret pairs

## Setup & deploying

- Build the Docker image
```shell
git clone https://github.com/streamcord/twproxy.git
docker build -t streamcord/twproxy:latest .
```
- Copy the `example.config.yml` file to `config.yml` and fill in the necessary values
- Create a container with port 8181 published
```shell
docker run -d --name twproxy --restart=always -p 8181:8181 streamcord/twproxy:latest
```

## Using with 3rd-party libraries

Twproxy uses a workaround with request headers from most 3rd-party libraries to support multiple API clients.

All you need to do is set the `ClientID` to the name of your service as defined in the config file, and set the `AppAccessToken` to the `auth` property for the service.

### Example

#### config.yml

(In your Twproxy config)

> **Warning**: Treat the value for `auth` like a password. If someone gets a hold of it, they can make requests to Twproxy on behalf of your API clients. You should also set the auth value to a very long random string.

```yaml
services:
  main:
    auth: "password123"
    client_id: "acfnw2z61xclweaajej39oxryiytvl"
    client_secret: "2cxntv598dtpbi08j73qh5nghz9kc1"
```

#### main.go

(In your client code)

```go
client, err := helix.NewClient(&helix.Options{
    ClientID:       "main",
    AppAccessToken: "password123",
    APIBaseURL:     "http://yourhost:8181/helix",
})
```

## Lint & format

```shell
gofmt -w -s .
golint ./...
```