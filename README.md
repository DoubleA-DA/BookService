# BookService

for running the application start the server and the client 

```go run server/s.go ```

```go run client/c.go```

for sending requests use ```http://localhost:50002/api/v1/```

sample datas given ```data.json review1.json review2.json```

example ```curl -X POST --data @data.json http://localhost:50002/api/v1/book```
