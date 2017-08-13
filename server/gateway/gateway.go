package gateway

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	API interface {
		GetID() int
		RegisterRoute(origin string, target string)
		GetSeed() int64
	}
	api struct {
		domain string
	}
)

func NewGateway(domain string) API {
	return &api{domain: domain}

}

func (api *api) GetID() int {
	body := mustHttpGet(fmt.Sprintf("%s/next_cluster_id", api.domain))
	id, err := strconv.Atoi(string(body))
	if err != nil {
		panic("could not init id")
	}
	return id
}

func (api *api) GetSeed() int64 {
	body := mustHttpGet(fmt.Sprintf("%s/seed", api.domain))
	seed, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		panic("can't init seed")
	}
	return seed
}

func (api *api) RegisterRoute(origin string, target string) {
	b, err := json.Marshal(struct {
		Origin string
		Target string
	}{Origin: origin, Target: target})

	if err != nil {
		panic(err)
	}

	mustHttpPost(fmt.Sprintf("%s/register_endpoint", api.domain), "application/json", b)
}
