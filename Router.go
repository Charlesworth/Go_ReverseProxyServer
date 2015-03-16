package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Stats struct {
	voteApiCalls int
}

var curentRecord = new(Stats)

func main() {

	router := httprouter.New()
	router.POST("/make", poo)
	router.GET("/stats", stats)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":299", nil))
}

func poo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	remote, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	curentRecord.voteApiCalls++

	fmt.Println("it works")

	proxy := httputil.NewSingleHostReverseProxy(remote)

	log.Println(r.URL)
	r.URL.Path = "/make"
	proxy.ServeHTTP(w, r)
}

func stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, curentRecord.voteApiCalls)
}
