# Go Chat App

Websocket based CloudEvents chat app.

## Usage

//TODO

## Running Locally

Test the frontend only:
```shel
cd web
npm start
```
Check on [React's README.md](./web/README.md) for more detailed info

Build the frontend:
```shell
cd web
npm run build4ko
```
or
```shell
cd web
npm run build
cp -r build/* ../../cmd/go-chat-client/kodata/www/
cd ..
```

Now to run the ko image locally:
```shell
KO_DATA_PATH=./cmd/go-chat-client/kodata go run cmd/go-chat-client/main.go
```
or

Run using ko and kubernetes:
Set your `KO_DOCKER_REPO` as stated [here](https://github.com/google/ko#choose-destination)
then run:
```shell
KO_DATA_PATH=./cmd/go-chat-client/kodata ko publish cmd/go-chat-client/main.go
```
and the run:
```shell
kubectl apply -f config/go-chat-client.yaml
```

### Cleanuo

```shell
ko delete -f config/go-chat-client.yaml
```

### Special Thanks

Thanks to @n3wscott for the original [https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye](https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye) from which I based most of this project's structure from.