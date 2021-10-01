package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	pf "BookService/protof"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)
type bookst struct{
	Id int64`json:"id"`
  	Name string`json:"name"`
  	Author []string`json:"author"`
  	Shortdesc string`json:"shortdesc"`
}
type review struct{
	Name string`json:"name"`
	Score int64`json:"score"`
	Text string`json:"text"`
	Id int64`json:"id"`
}
var(
	clicon,err=grpc.Dial("localhost:50001",grpc.WithInsecure(),grpc.WithBlock())
	cli=pf.NewSRVClient(clicon)
)
func getbook(w http.ResponseWriter,r*http.Request){
	w.Header().Set("Content-Type","application/json")
	vars:=mux.Vars(r)
	v,_:=strconv.Atoi(vars["limit"])
	snd:=&pf.Lim{L:int64(v)}
	if rcv,err1:=cli.GetBook(r.Context(),snd);err1==nil{
		json.NewEncoder(w).Encode(rcv.Ak)
		for _,it:=range rcv.BookList{
			json.NewEncoder(w).Encode(it)
		}
	}
}
func insbook(w http.ResponseWriter,r*http.Request){
	w.Header().Set("Content-Type","application/json")
	//x,_:=ioutil.ReadAll(r.Body)
	//fmt.Println(string(x))
	var bk bookst
	json.NewDecoder(r.Body).Decode(&bk)
	fmt.Println(bk)
	snd:=&pf.Book{BookId:bk.Id,Name:bk.Name,Author:bk.Author,ShortDesc:bk.Shortdesc}
	if rcv,err1:=cli.InsBook(r.Context(),snd);err1==nil{
		json.NewEncoder(w).Encode(rcv)
	}
}
func insreview(w http.ResponseWriter,r*http.Request){
	w.Header().Set("Content-Type","application/json")
	var rv review
	json.NewDecoder(r.Body).Decode(&rv)
	fmt.Println(rv)
	snd:=&pf.Review{BookId:rv.Id,Name:rv.Name,Score:rv.Score,Text:rv.Text}
	if rcv,err:=cli.InsReview(r.Context(),snd);err==nil{
		json.NewEncoder(w).Encode(rcv)
	}
}
func getreview(w http.ResponseWriter,r*http.Request){
	w.Header().Set("Content-Type","application/json")
	vars:=mux.Vars(r)
	v,_:=strconv.Atoi(vars["bookid"])
	snd:=&pf.Rv{BookId:int64(v)}
	if rcv,err1:=cli.GetReview(r.Context(),snd);err1==nil{
		json.NewEncoder(w).Encode(rcv.Ak)
		for _,it:=range rcv.ReviewList{
			json.NewEncoder(w).Encode(it)
		}
	}
}
/*
curl -X POST --data @data.json http://localhost:50002/api/v1/book
curl -X GET http://localhost:50002/api/v1/review/1
curl -d '{"id":1,"name":"xyz","author":["fg","h"],"shortdesc":"sa"}' -H "Content-Type: application/json" -X POST http://localhost:50002/api/v1/book
*/
func main(){
	if err!=nil{
		panic("can't connect")
	}
	mx:=mux.NewRouter()
	mx.HandleFunc("/api/v1/book/{limit}",getbook).Methods("GET")
	mx.HandleFunc("/api/v1/book",insbook).Methods("POST")
	mx.HandleFunc("/api/v1/review",insreview).Methods("POST")
	mx.HandleFunc("/api/v1/review/{bookid}",getreview).Methods("GET")
	if err1:=http.ListenAndServe(":50002",mx);err1!=nil{
		panic("Error can't connect to http")
	}
}