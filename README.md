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
cp -r build/* ../../cmd/knative-go-chat-app/kodata/www/
cd ..
```

Now to run the ko image locally:
```shell
KO_DATA_PATH=./cmd/knative-go-chat-app/kodata go run cmd/knative-go-chat-app/main.go
```

## Running on a Cluster

Run using ko and kubernetes:
Set your `KO_DOCKER_REPO` as stated [here](https://github.com/google/ko#choose-destination)
then run:
simply run:
```shell
ko apply -f config/knative-go-chat-app.yaml
```
To build your own image do:
```shell
KO_DATA_PATH=./cmd/knative-go-chat-app/kodata ko publish -B ./cmd/knative-go-chat-app/main.go
```
Then modify the `image URI` in the `config/knative-go-chat-app.yaml` file 
and finally run:
```shell
kubectl apply -f config/knative-go-chat-app.yaml
```
### Cleanup

```shell
ko delete -f config/knative-go-chat-app.yaml
```

### Special Thanks

Thanks to @n3wscott for the original [https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye](https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye) from which I based most of this project's structure from.
