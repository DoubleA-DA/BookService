package db

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
)


type CB struct{
	Collection_book *gocb.Collection
	Collection_review *gocb.Collection
	Collection_infos *gocb.Collection
	Scope *gocb.Scope
	Cluster *gocb.Cluster
	Infos *gocb.LookupInResult
}
type bookst struct{
  	Name string`json:"name,omitempty"`
  	Author []string`json:"author,omitempty"`
  	Shortdesc string`json:"shortdesc,omitempty"`
	Counter int`json:"counter,omitempty"`
}
type review struct{
	Name string`json:"name"`
	Score int64`json:"score"`
	Text string`json:"text"`

}
type blank struct{
	Review[]int`json:"review"`
}
type info struct{
	Last int`json:"last"`
	Total int`json:"total,omitempty"`
}
var(
	bucket_name="bookdb"
	scope_name="book_details"
	collection_books="books"
	collection_reviews="reviews"
	collection_infos="infos"
	max_store int64
	cbh CB
)

func connectDB()(*gocb.Cluster){
	cluster,err:=gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username:"DA",
			Password:"Alcatraz100",
		})
	if err!=nil{
		panic("can't connect to DB")
	}
	cbh.Cluster=cluster
	return cluster
}
func Initializer()(bool){
	cluster:=connectDB()
	var bucket *gocb.Bucket
	bucmgr:=cluster.Buckets()
	if _,err:=bucmgr.GetBucket(bucket_name,nil);err!=nil{
	if err:=bucmgr.CreateBucket(gocb.CreateBucketSettings{
		BucketSettings:gocb.BucketSettings{
			Name:bucket_name,
			FlushEnabled:false,
			ReplicaIndexDisabled:true,
			RAMQuotaMB:1000,
			NumReplicas:0,
			BucketType:gocb.CouchbaseBucketType,
		},
		ConflictResolutionType:gocb.ConflictResolutionTypeSequenceNumber,
	},nil);err!=nil{
		panic("Can't create bucket")
	}}
	bucket=cluster.Bucket(bucket_name)
	bmgr:=bucket.Collections()
	s,_:=bmgr.GetAllScopes(nil)
	if s[0].Name!=scope_name{
		if err:=bmgr.CreateScope(scope_name,nil);err!=nil{
			panic("can't create scope")
		}
	}
	check:=func()bool{
		for _,it:=range s[0].Collections{
			if it.Name==collection_books{
				return false
			}
		}
		return true
		}
	if check(){
		if err:=bmgr.CreateCollection(
			gocb.CollectionSpec{
				Name:collection_books,
				ScopeName:scope_name,
			},nil);err!=nil{
				panic("can't create collection")
			}
	}
	check=func()bool{
		for _,it:=range s[0].Collections{
			if it.Name==collection_reviews{
				return false
			}
		}
		return true
		}
	if check(){
		if err:=bmgr.CreateCollection(
			gocb.CollectionSpec{
				Name:collection_reviews,
				ScopeName:scope_name,
			},nil);err!=nil{
				panic("can't create collection")
			}
	}
	check=func()bool{
		for _,it:=range s[0].Collections{
			if it.Name==collection_infos{
				return false
			}
		}
		return true
		}
	if check(){
		if err:=bmgr.CreateCollection(
			gocb.CollectionSpec{
				Name:collection_infos,
				ScopeName:scope_name,
			},nil);err!=nil{
				panic("can't create collection")
			}
	}
	scope:=bucket.Scope(scope_name)
	collection:=scope.Collection(collection_books)
	cbh.Collection_book=collection
	collection=scope.Collection(collection_reviews)
	cbh.Collection_review=collection
	collection=scope.Collection(collection_infos)
	cbh.Collection_infos=collection
	cbh.Scope=scope
	ops:=[]gocb.LookupInSpec{
		gocb.GetSpec("last",nil),
	}
	_,err:=cbh.Collection_infos.LookupIn("info",ops,nil)
	if err!=nil{
		var in info
		in.Last=1;
		if _,err=cbh.Collection_infos.Insert("info",&in,&gocb.InsertOptions{});err!=nil{
			return false
		}
	}
	cbh.Infos,_=cbh.Collection_infos.LookupIn("info",ops,nil)
	fmt.Println("Enter max storage limit")
	fmt.Scanf("%d",&max_store)
	return true
}
func Addbook(id string,name string,author[]string,shortdesc string)(string){
	var bk bookst
	var bl blank
	bl.Review=[]int{}
	var in info

	cbh.Infos.ContentAt(0,&in.Last)//Last 
	fmt.Println(in)
	ops:=[]gocb.LookupInSpec{//checking counter
		gocb.GetSpec("counter",nil),
	}
	_,err:=cbh.Collection_review.Insert(id,&bl,&gocb.InsertOptions{})
	if err!=nil{
		return "Can't insert bookinfo, check index"
	}

	res,err:=cbh.Collection_book.LookupIn(fmt.Sprintf("%d",in.Last),ops,nil)
	if err!=nil{
		bk.Counter=1
		cbh.Collection_book.Insert(fmt.Sprintf("%d",in.Last),&bk,nil)
		bk.Counter=0
		bk.Name=name
		bk.Author=author
		bk.Shortdesc=shortdesc
		mops := []gocb.MutateInSpec{
			gocb.InsertSpec(id, &bk, &gocb.InsertSpecOptions{}),
		}
		cbh.Collection_book.MutateIn(fmt.Sprintf("%d",in.Last), mops, &gocb.MutateInOptions{})
	}else{

		res.ContentAt(0,&bk.Counter)
		if(bk.Counter<int(max_store)){
			bk.Name=name
			bk.Author=author
			bk.Shortdesc=shortdesc

			ops=[]gocb.LookupInSpec{//checking hash
				gocb.ExistsSpec(id,nil),
			}
			res,_:=cbh.Collection_book.LookupIn(fmt.Sprintf("%d",in.Last),ops,nil)
			if res.Exists(0){
				return "Can't insert bookinfo, check index"
			}

			ops1:= []gocb.MutateInSpec{//increasing counter
				gocb.IncrementSpec("counter", 1, &gocb.CounterSpecOptions{}),
			}
			_, err= cbh.Collection_book.MutateIn(fmt.Sprintf("%d",in.Last), ops1, nil)
			if err != nil {
				panic(err)
			}
			bk.Counter=0
			mops := []gocb.MutateInSpec{//inserting book
				gocb.InsertSpec(id, &bk, &gocb.InsertSpecOptions{}),
			}
			_, err = cbh.Collection_book.MutateIn(fmt.Sprintf("%d",in.Last), mops, &gocb.MutateInOptions{})
			if err != nil {
				panic(err)
			}
			}else{
			ops1:= []gocb.MutateInSpec{//increasing last
				gocb.IncrementSpec("last", 1, &gocb.CounterSpecOptions{}),
			}
			_, err= cbh.Collection_infos.MutateIn("info", ops1, nil)
			if err != nil {
				panic(err)
			}
			in.Last+=1
			bk.Counter=1
			cbh.Collection_book.Insert(fmt.Sprintf("%d",in.Last),&bk,nil)
			bk.Counter=0
			bk.Name=name
			bk.Author=author
			bk.Shortdesc=shortdesc
			mops := []gocb.MutateInSpec{//Inserting book
				gocb.InsertSpec(id, &bk, &gocb.InsertSpecOptions{}),
			}
			cbh.Collection_book.MutateIn(fmt.Sprintf("%d",in.Last), mops, &gocb.MutateInOptions{})
			}
	}

	
	ops=[]gocb.LookupInSpec{//checking if "1" exists or not
		gocb.ExistsSpec(fmt.Sprintf("%d",in.Last),nil),
	}
	res,_=cbh.Collection_infos.LookupIn("info",ops,nil)

	if !res.Exists(0){
		in.Total=1

		mops1 := []gocb.MutateInSpec{//inserting "1"
			gocb.InsertSpec(fmt.Sprintf("%d",in.Last), &in.Total, &gocb.InsertSpecOptions{}),
		}
		_, err = cbh.Collection_infos.MutateIn("info", mops1, &gocb.MutateInOptions{})
		
		if err != nil {
			panic(err)
		}
	}else{
		ops1:= []gocb.MutateInSpec{//increasing total
			gocb.IncrementSpec(fmt.Sprintf("%d",in.Last), 1, &gocb.CounterSpecOptions{}),
		}
		_, err = cbh.Collection_infos.MutateIn("info", ops1, nil)
		if err != nil {
			panic(err)
		}
	}

	if err!=nil{
		return "Can't insert bookinfo, check index"
	}
	return "ok"
}
func Addreview(ind string,name string,score int64,text string)(string){
	var rv review
	rv.Name=name
	rv.Score=score
	rv.Text=text

	mut:=[]gocb.MutateInSpec{
		gocb.ArrayAppendSpec("review",rv,nil),
	}
	

	_,err:=cbh.Collection_review.MutateIn(ind,mut,&gocb.MutateInOptions{
	})

	if err!=nil{
		return "check index"
	}
	return "ok"
	
} 
func Retbook(lim int64)([]bookst,bool){
	var bk []bookst
	var bk1 bookst

	n:=lim/max_store
	rem:=lim%max_store
	var i int64
	var dat map[string]map[string]string//name,shortdesc
	var dat1 map[string]map[string][]string//author
	for i=1;i<=n;i++{
		res,err:=cbh.Collection_book.Get(fmt.Sprintf("%d",i),nil)
		if err!=nil{
			rem=0
			break
		}
		res.Content(&dat)
		res.Content(&dat1)
	}
	for k:=range dat{
		bk1.Name=""
		for i,j:=range dat[k]{
			if(len(j)==0){continue}
			if(i=="name"){
				bk1.Name=j
			}
			if(i=="shortdesc"){
				bk1.Shortdesc=j
			}
		}
		for _,j:=range dat1[k]{
			if(j==nil){continue}
			bk1.Author=j
		}
		if len(bk1.Name)==0{continue}
		bk=append(bk,bk1)
	}
	for k:=range dat{
		delete(dat,k)
		delete(dat1,k)
	}
	if rem!=0{
		res,err:=cbh.Collection_book.Get(fmt.Sprintf("%d",i),nil)
		if err!=nil{
			return bk,false
		}
		res.Content(&dat)
		res.Content(&dat1)
	
		for k:=range dat{
			bk1.Name=""
			if rem==0{
				break
			}
			for i,j:=range dat[k]{
				if(len(j)==0){continue}
				if(i=="name"){
					bk1.Name=j
				}
				if(i=="shortdesc"){
					bk1.Shortdesc=j
				}
			}
			for _,j:=range dat1[k]{
				if(j==nil){continue}
				bk1.Author=j
			}
			if len(bk1.Name)==0{continue}
			rem--
			bk=append(bk,bk1)
		}
	}
	if len(bk1.Name)==0{
		return bk,false
	}
	return bk,true
}
func Retreview(id string)([]review,bool){
	var rv []review
	ops:=[]gocb.LookupInSpec{
		gocb.GetSpec("review",nil),
	}
	res,err:=cbh.Collection_review.LookupIn(id,ops,nil)
	if err!=nil{
		return rv,false
	}
	err=res.ContentAt(0,&rv)
	if err!=nil{
		return rv,false
	}
	fmt.Println(rv)
	if rv[0].Score<0{
		return rv,false
	}
	return rv,true
}