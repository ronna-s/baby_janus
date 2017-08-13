package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/wwgberlin/baby_janus/api/cluster"
)

type (
	Endpoint struct {
		Origin string
		Target string
	}
)

/*
	redirectHandler returns a handler to redirect the request to
 */

func redirectHandler(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusFound)
	}
}

/*
	registerEndpoint handles requests to registers routes origin - target
 */
func registerEndpoint(response http.ResponseWriter, request *http.Request) {
	var endpoint Endpoint
	if request != nil && request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &endpoint); err != nil {
			response.WriteHeader(http.StatusBadRequest)
		}
		http.HandleFunc(endpoint.Origin, redirectHandler(endpoint.Target))
		response.WriteHeader(http.StatusCreated)
	}
}

/*
	incrClusterId - returns handler to increment the cluster servers size
 */

func incrClusterId(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.IncrClusterId()))
	}
}

/*
	helloUser fetches the parts from all the APIs registered to the cluster
 */
func helloUser(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodies := []string{}
		for _, path := range GetSlices() {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080%s", path))
			if err != nil {
				panic(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			fmt.Println(body)
			bodies = append(bodies, string(body))
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, strings.Join(bodies, ""))
	}
}

func getSeed(c cluster.Cluster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%v", c.GetSeed()))
	}
}

func main() {
	c := cluster.NewCluster()

	/*
		register your initial routes for the API here
	 */
	http.HandleFunc("/", helloUser(c))
	http.HandleFunc("/next_cluster_id", incrClusterId(c))
	http.HandleFunc("/seed", getSeed(c))
	http.HandleFunc("/register_endpoint", registerEndpoint)

	http.ListenAndServe(":8080", nil)

}

func GetSlices() []string {
	res := make([]string, c.numSlices)
	for i := range res {
		res[i] = fmt.Sprintf("%v", c.slicer(i))
	}
	return res

}
