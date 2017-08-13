package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"math/rand"

	"github.com/wwgberlin/baby_janus/server/cluster"
	"github.com/wwgberlin/baby_janus/server/gateway"
)

func main() {
	<-time.After(10 * time.Second) //give the server some time to start
	myDomain := os.Getenv("HOSTNAME")
	api := gateway.NewGateway("http://baby_janus_api:8080")
	registerRoutes(api, api.GetID(), myDomain)
	rand.Seed(api.GetSeed())
	http.ListenAndServe(":8080", nil)

}

func registerRoutes(api gateway.API, id int, myDomain string) {
	routes := getRoutes(id)
	for i := range routes {
		route := routes[i]
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(fmt.Sprintf(".%s", route))
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, string(b))
		})
		api.RegisterRoute(fmt.Sprintf(route), fmt.Sprintf("http://%s:8080%s", myDomain, route))
	}

}

func getRoutes(id int) (slices []string) {
	c := cluster.NewCluster()
	slices = c.GetInstanceSlices(id)
	return
}
