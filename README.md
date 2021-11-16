# BookService

for running the application start the server and the client 

```go run server/s.go ```

```go run client/c.go```

for the limit of books per doc, enter the limit in the grpc server when prompted

for sending requests use ```http://localhost:50002/api/v1/```

sample datas given 

example ```curl -X POST --data @data.json http://localhost:50002/api/v1/book```
