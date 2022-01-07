# KeepIt



```
vscode-proto3


sudo apt install -y protobuf-compiler
go mod init github.com/SmsS4/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=.  --go-grpc_opt=paths=source_relative route_guide.proto
protoc --go_out=plugins=grpc:chat srvc.proto 

```
