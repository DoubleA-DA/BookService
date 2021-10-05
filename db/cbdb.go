package db

import (
	"fmt"
	//"sync"
	//"time"
	"github.com/couchbase/gocb/v2"

	//"github.com/go-playground/validator/v10"
)


type CB struct{
	Collection *gocb.Collection
	Scope *gocb.Scope
	//lock sync.Mutex
}
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
}
var(
	bucket_name="bookdb"
	scope_name="books_det"
	collection_name="infos"
	primary_key="#primary"
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
	if s[0].Collections[0].Name!=collection_name{
		if err:=bmgr.CreateCollection(
			gocb.CollectionSpec{
				Name:collection_name,
				ScopeName:scope_name,
			},nil);err!=nil{
				panic("can't create collection")
			}
	}
	scope:=bucket.Scope(scope_name)
	collection:=scope.Collection(collection_name)
	cbh.Collection=collection
	cbh.Scope=scope
	return true
}
func Addbook(id int64,name string,author[]string,shortdesc string)(string){
	var bk bookst
	bk.Id=id
	bk.Name=name
	bk.Author=author
	bk.Shortdesc=shortdesc
	bk.Review=[]int{}
	/*validate:=validator.New()
	if err:=validate.Struct(bk);err!=nil{
		v:=err.(validator.ValidationErrors)
		return fmt.Sprintf("%s",v)
	}*/
	//cbh.lock.Lock()
	_,err:=cbh.Collection.Insert(fmt.Sprintf("%d",bk.Id),&bk,&gocb.InsertOptions{})
	//cbh.lock.Unlock()
	if err!=nil{
		return "Can't insert bookinfo, check index"
	}
	return "ok"
}
func Addreview(ind int64,name string,score int64,text string)(string){
	var rv review
	rv.Name=name
	rv.Score=score
	rv.Text=text
	/*validate:=validator.New()
	if err:=validate.Struct(rv);err!=nil{
		v:=err.(validator.ValidationErrors)
		return fmt.Sprintf("%s",v)
	}*/
	mut:=[]gocb.MutateInSpec{
		gocb.ArrayAppendSpec("review",rv,nil),
	}
	//cbh.lock.Lock()
	/*lck,err:= cbh.Collection.Get(fmt.Sprintf("%d",ind), nil)
	if err!=nil{
		panic("can't fetch")
	}*/
	//lck1:= lck.Cas()
	_,err:=cbh.Collection.MutateIn(fmt.Sprintf("%d",ind),mut,&gocb.MutateInOptions{
	})
	//cbh.Collection.Unlock(fmt.Sprintf("%d",ind), lck1, nil)
	//cbh.lock.Unlock()
	if err!=nil{
		return "check index"
	}
	return "ok"
	
} 
func Retbook(lim int64)([]bookst,bool){
	var bk []bookst
	query:="select bucket_id from system:indexes where name=$1";
	rows,err:=cbh.Scope.Query(query,&gocb.QueryOptions{PositionalParameters:[]interface{}{primary_key}})
	if err!=nil{
		return bk,false
	}
	var s map[string]string
	for rows.Next(){
		rows.Row(&s)
	}
	if s["bucket_id"]!=bucket_name{
		query="create primary index on `bookdb`.books_det.infos";
		_,err=cbh.Scope.Query(query,nil)
		if err!=nil{
			panic("can't create index")
		}
	}
	query="select name,id,author,shortdesc from `bookdb`.books_det.infos limit $1"
	rows,err=cbh.Scope.Query(query,&gocb.QueryOptions{PositionalParameters:[]interface{}{lim}})
	if err!=nil{
		return bk,false
	}
	var v bookst
	for rows.Next(){
		rows.Row(&v)
		bk=append(bk,v)
		fmt.Println(bk)
	}
	if len(v.Name)==0{return bk,false}
	return bk,true
}
func Retreview(id int64)([]review,bool){
	var rv []review
	ops:=[]gocb.LookupInSpec{
		gocb.GetSpec("review",nil),
	}
	res,err:=cbh.Collection.LookupIn(fmt.Sprintf("%d",id),ops,nil)
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