# mexc-golang-sdk

## To update websocket protos
 - init and update git submodule
 - add to proto file line `option go_package = "websocket/dto";`
 - run `make proto`