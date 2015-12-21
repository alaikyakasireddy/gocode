package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"github.com/julienschmidt/httprouter"

// 	"golang.org/x/net/context"
// 	"golang.org/x/oauth2/google"

// 	"google.golang.org/appengine"
// 	"google.golang.org/cloud"
// 	"google.golang.org/cloud/storage"
// 	"io/ioutil"
// )

// const gcsBucket = "collabo-write-stories"
// const aeId = "collabo-write"

// /*
// func init() {
// 	http.HandleFunc("/put", handlePut)
// 	http.HandleFunc("/get", handleGet)
// 	http.HandleFunc("/list", handleList)
// }
// */

// func getCloudContext(req *http.Request) (context.Context, error) {
// 	jsonKey, err := ioutil.ReadFile("Collabowrite-5737485e8364.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	conf, err := google.JWTConfigFromJSON(
// 		jsonKey,
// 		storage.ScopeFullControl,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx := appengine.NewContext(req)
// 	hc := conf.Client(ctx)
// 	return cloud.NewContext(aeId, hc), nil
// }

// func handlePut(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
// 	cctx, err := getCloudContext(req)
// 	if err != nil {
// 		http.Error(res, "ERROR GETTING CCTX: "+err.Error(), 500)
// 		return
// 	}
	
// 	fileName := "test.txt"
// 	contentToWrite := "hello world!"

// 	writer := storage.NewWriter(cctx, gcsBucket, fileName)
// 	io.WriteString(writer, contentToWrite)
// 	err = writer.Close()
// 	if err != nil {
// 		http.Error(res, "ERROR WRITING TO BUCKET: "+err.Error(), 500)
// 		return
// 	}
// }

// func handleGet(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
// 	cctx, err := getCloudContext(req)
// 	if err != nil {
// 		http.Error(res, "ERROR GETTING CCTX: "+err.Error(), 500)
// 		return
// 	}
	
// 	fileName := "test.txt"

// 	rdr, err := storage.NewReader(cctx, gcsBucket, fileName)
// 	if err != nil {
// 		http.Error(res, err.Error(), 500)
// 		return
// 	}
// 	defer rdr.Close()

// 	io.Copy(res, rdr)
// }

// func handleList(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
// 	cctx, err := getCloudContext(req)
// 	if err != nil {
// 		http.Error(res, "ERROR GETTING CCTX: "+err.Error(), 500)
// 		return
// 	}

// 	var query *storage.Query
// 	objs, err := storage.ListObjects(cctx, gcsBucket, query)
// 	if err != nil {
// 		http.Error(res, err.Error(), 500)
// 		return
// 	}
// 	for _, obj := range objs.Results {
// 		fmt.Fprintln(res, obj.Name)
// 	}
// }
