package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	//"io/ioutil"
	pf "BookService/protof"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type bookst struct{
	Id int64`json:"id" validate:"required"`
  	Name string`json:"name" validate:"required"`
  	Author []string`json:"author" validate:"required"`
  	Shortdesc string`json:"shortdesc"`
	Review []int`json:"review"`
}
type review struct{
	Name string`json:"name" validate:"required"`
	Score int64`json:"score" validate:"required,gte=1,lte=5"`
	Text string`json:"text"`
	Id int64`json:"id"`
}

var (
	clicon, err = grpc.Dial("localhost:50001", grpc.WithInsecure(), grpc.WithBlock())
	cli         = pf.NewSRVClient(clicon)
)

func getbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	v, _ := strconv.Atoi(vars["limit"])
	snd := &pf.Limit{L: int64(v)}
	if rcv, err1 := cli.GetBook(r.Context(), snd); err1 == nil {
		json.NewEncoder(w).Encode(rcv.Ak)
		for _, it := range rcv.BookList {
			json.NewEncoder(w).Encode(it)
		}
	}
}
func insbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//x,_:=ioutil.ReadAll(r.Body)
	//fmt.Println(string(x))
	var bk bookst
	json.NewDecoder(r.Body).Decode(&bk)
	fmt.Println(bk)
	validate:=validator.New()
	if err:=validate.Struct(bk);err!=nil{
		v:=err.(validator.ValidationErrors)
		json.NewEncoder(w).Encode(fmt.Sprintf("%s",v))
	}else{
		snd := &pf.Book{BookId: bk.Id, Name: bk.Name, Author: bk.Author, ShortDesc: bk.Shortdesc}
		if rcv, err1 := cli.InsBook(r.Context(), snd); err1 == nil {
			json.NewEncoder(w).Encode(rcv)
		}
	}
}
func insreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rv review
	json.NewDecoder(r.Body).Decode(&rv)
	fmt.Println(rv)
	validate:=validator.New()
	if err:=validate.Struct(rv);err!=nil{
		v:=err.(validator.ValidationErrors)
		json.NewEncoder(w).Encode(fmt.Sprintf("%s",v))
	}else{
		snd := &pf.Review{BookId: rv.Id, Name: rv.Name, Score: rv.Score, Text: rv.Text}
		if rcv, err := cli.InsReview(r.Context(), snd); err == nil {
			json.NewEncoder(w).Encode(rcv)
		}
	}
}
func getreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	v, _ := strconv.Atoi(vars["bookid"])
	snd := &pf.Review{BookId: int64(v)}
	if rcv, err1 := cli.GetReview(r.Context(), snd); err1 == nil {
		json.NewEncoder(w).Encode(rcv.Ak)
		for _, it := range rcv.ReviewList {
			json.NewEncoder(w).Encode(it)
		}
	}
}

/*
curl -X POST --data @data.json http://localhost:50002/api/v1/book
curl -X GET http://localhost:50002/api/v1/review/1
curl -d '{"id":1,"name":"xyz","author":["fg","h"],"shortdesc":"sa"}' -H "Content-Type: application/json" -X POST http://localhost:50002/api/v1/book
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative prt.proto
*/
func main() {
	if err != nil {
		panic("can't connect")
	}
	mx := mux.NewRouter()
	mx.HandleFunc("/api/v1/book/{limit}", getbook).Methods("GET")
	mx.HandleFunc("/api/v1/book", insbook).Methods("POST")
	mx.HandleFunc("/api/v1/review", insreview).Methods("POST")
	mx.HandleFunc("/api/v1/review/{bookid}", getreview).Methods("GET")
	if err1 := http.ListenAndServe(":50002", mx); err1 != nil {
		panic("Error can't connect to http")
	}
}
