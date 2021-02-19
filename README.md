# Tasks
Basic Golang gRPC service that implements routes for toDoList

## Proto Compile
To compile gRPC service run:

``
protoc -I $GOPATH/src --go_out=$GOPATH/src $GOPATH/src/tasks/internal/proto-files/messages/tasks.proto
``

and 

``
protoc -I $GOPATH/src --go_out=plugins=grpc:$GOPATH/src $GOPATH/src/tasks/internal/proto-files/service/tasks-service.proto
``