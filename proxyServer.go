package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

var proxyCalls int

func main() {

	router := httprouter.New()
	router.GET("/test1", reverseProxy)
	router.GET("/test2", serialCompositionProxy)
	router.GET("/test3", parrallelCompositionProxy)
	router.GET("/stats", stats)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	//Instead of using http.listen... you can create a custom server to create custom behavior
	//such as read and write timeouts or max header bytes:
	// s := &http.Server{
	// 	Addr:           ":8080",
	// 	Handler:        router,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// log.Fatal(s.ListenAndServe())
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

func serialCompositionProxy(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body := make([]string, 4)

	for i := 0; i < 4; i++ {
		resp, err := http.Get("http://localhost:3000")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		str, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		body[i] = string(str)
	}

	fmt.Fprintln(res, body)
}

func parrallelCompositionProxy(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body := make([]string, 4)
	ret := make(chan string)

	for i := 0; i < 4; i++ {
		go func() {
			resp, err := http.Get("http://localhost:3000")
			if err != nil {
				panic(err)
			}

			str, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			resp.Body.Close()
			ret <- string(str)
		}()
	}

	for returns := 0; returns < 4; {
		body[returns] = <-ret
		returns++
	}
	fmt.Fprintln(res, body)
}

func stats(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(res, proxyCalls)
}
