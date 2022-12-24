package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func handler(res http.ResponseWriter, req *http.Request) {
	// Example of parsing GET or POST Query Params

	log.Println("----------------------------------")
	log.Println("RequestURI", req.RequestURI)
	log.Println("RemoteAddr", req.RemoteAddr)
	log.Println("req.Host", req.Host)
	log.Println("req.Method", req.Method)
	log.Println("req.URL.Path", req.URL.Path)
	//log.Println("URL", r.URL)
	log.Println("----------------------------------")

	//url := "http://mockbin.com/request?foo=bar&foo=baz"
	url := "http://localhost:6000"

	payload := strings.NewReader("{\"foo\": \"bar\"}")

	req1, _ := http.NewRequest(req.Method, url, payload)

	req1.Header.Add("cookie", "foo=bar; bar=baz")
	req1.Header.Add("accept", "application/json")
	req1.Header.Add("content-type", "application/json")

	res1, _ := http.DefaultClient.Do(req1)

	defer res1.Body.Close()
	body, _ := ioutil.ReadAll(res1.Body)

	fmt.Println(res1)
	fmt.Println(string(body))

	res.WriteHeader(res1.StatusCode)
	res.Write(body)

}

func main() {
	fmt.Println("WARNING: This is an example, but not really safe.")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8001", Log(http.DefaultServeMux))
}
