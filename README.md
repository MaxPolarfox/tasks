# Tasks
Basic Golang gRPC service which implements routes for toDoList

## Proto Compile
To compile gRPC service run:

``
protoc -I $GOPATH/src/github.com/MaxPolarfox --go_out=$GOPATH/src/github.com/MaxPolarfox $GOPATH/src/github.com/MaxPolarfox/tasks/internal/proto-files/messages/tasks.proto
``

and 

``
protoc -I $GOPATH/src/github.com/MaxPolarfox --go_out=plugins=grpc:$GOPATH/src/github.com/MaxPolarfox $GOPATH/src/github.com/MaxPolarfox/tasks/internal/proto-files/service/tasks-service.proto
``