package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type (
	server struct {
		id     string
		routes []string
	}
)

func getId() string {
	resp, err := http.Get("http://baby_janus_api:8080/next_cluster_id")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(body)
}

func getRoutes(id string) (slices []string) {
	body := mustHttpPost("http://baby_janus_api:8080/get_instance_slices", "text/plain", []byte(id))
	if err := json.Unmarshal(body, &slices); err != nil {
		panic(err)
	}
	return
}

func registerRoute(origin string, target string) {
	b, err := json.Marshal(struct {
		Origin string
		Target string
	}{Origin: origin, Target: target})

	if err != nil {
		panic(err)
	}

	mustHttpPost(fmt.Sprintf("http://baby_janus_api:8080/register_endpoint"), "application/json", b)
}

func main() {
	<-time.After(10 * time.Second) //give the server some time to start
	myDomain := os.Getenv("HOSTNAME")
	id := getId()
	routes := getRoutes(id)
	//server{id: getId(), routes: routes}
	for i := range routes {
		route := routes[i]
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(fmt.Sprintf(".%s", route))
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, string(b))
		})
		registerRoute(fmt.Sprintf(route), fmt.Sprintf("http://%s:8080%s", myDomain, route))
	}
	http.ListenAndServe(":8080", nil)

}

func mustHttpPost(url string, contentType string, body []byte) []byte {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}
