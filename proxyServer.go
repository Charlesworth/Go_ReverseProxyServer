package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

var proxyCalls int

func main() {

	router := httprouter.New()
	router.GET("/proxy", reverseProxy)
	router.GET("/stats", stats)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func reverseProxy(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	proxyTargetURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	proxyCalls++

	proxy := httputil.NewSingleHostReverseProxy(proxyTargetURL)

	log.Println(req.URL)
	req.URL.Path = "/"
	proxy.ServeHTTP(res, req)
}

func stats(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(res, proxyCalls)
}
