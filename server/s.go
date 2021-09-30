package main

import (
	"context"
	"fmt"
	"net"
	"proj/db"
	pf "proj/protof"
	"google.golang.org/grpc"
	"net/http"
)
type servr struct{
	pf.UnimplementedSRVServer
}

func(s*servr)InsBook(ctx context.Context,req*pf.Book)(*pf.Ack,error){
	if v:=db.Addbk(req.GetBookId(),req.GetName(),req.GetAuthor(),req.GetShortDesc());v!="ok"{
		return &pf.Ack{Status:http.StatusBadRequest,Msg:v},nil
	}
	return &pf.Ack{Status:http.StatusCreated,Msg:"Success"},nil
}

func(s*servr)InsReview(ctx context.Context,req*pf.Review)(*pf.Ack,error){
	if v:=db.Addrv(req.GetBookId(),req.GetName(),req.GetScore(),req.GetText());v!="ok"{
		return &pf.Ack{Status:http.StatusBadRequest,Msg:v},nil
	}
	return &pf.Ack{Status:http.StatusCreated,Msg:"Success"},nil
}

func(s*servr)GetBook(ctx context.Context,req*pf.Lim)(*pf.Dispbook,error){
	var s1 []*pf.Book
	if bk,err:=db.Retbk(req.GetL());err{
		for _,b:=range bk{
			s1=append(s1,&pf.Book{Name:b.Name,Author:b.Author,ShortDesc:b.Shortdesc,BookId:b.Id})
		}
	}else{return &pf.Dispbook{Ak:&pf.Ack{Status:http.StatusNotFound,Msg:"No Books entered"}},nil}
	return &pf.Dispbook{Ak:&pf.Ack{Status:http.StatusFound},BookList:s1},nil
}

func(s*servr)GetReview(ctx context.Context,req*pf.Rv)(*pf.Dispreview,error){
	var s1 []*pf.Review
	if rv,err:=db.Retrv(req.GetBookId());err{
		for i:=0;i<len(rv);i++{
			s1=append(s1,&pf.Review{Name:rv[i].Name,Score:rv[i].Score,Text:rv[i].Text})
		}
	}else{return &pf.Dispreview{Ak:&pf.Ack{Status:http.StatusNotFound,Msg:"No data"}},nil}
	return &pf.Dispreview{Ak:&pf.Ack{Status:http.StatusFound},ReviewList:s1},nil
}



func main(){
	if db.Initializer(){
		fmt.Println("connected to DB")
	}
	lis,err1:=net.Listen("tcp",":50001")
	if err1!=nil{
		panic("can't connect to server")
	}
	srvr:=grpc.NewServer()
	pf.RegisterSRVServer(srvr,&servr{})
	if err:=srvr.Serve(lis);err!=nil{
		panic("can't serve")
	}
	
}