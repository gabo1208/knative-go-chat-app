# Go Chat Client

Websocket based CloudEvents chat viewer.

## Usage

//TODO

## Running Locally

Test the frontend only:
```shel
cd pkg/frontend
npm start
```
Check on [React's README.md](./pkg/frontend/README.md) for more detailed info

Build the frontend:
```shell
cd pkg/frontend
npm run build4ko
```
or
```shell
cd pkg/frontend
npm run build
cp -r build/* ../../cmd/go-chat-client/kodata/www/
cd ../..
```

Now to run the ko image locally:
```shell
KO_DATA_PATH=./cmd/go-chat-client/kodata go run cmd/go-chat-client/main.go
```

### From Source

```shell
ko apply -f config/go-chat-client.yaml
```

### Special Thanks

Thanks to @n3wscott for the original [https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye](https://github.com/n3wscott/sockeye/tree/main/cmd/sockeye) from which I based most of this project's structure from.