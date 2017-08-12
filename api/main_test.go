package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"runtime"
)

type (
	clusterMock struct {
		iterator func() int
	}
)

func TestRedirect(t *testing.T) {
	startServer()

	calledBack := false
	origin := "/origin"
	target := "http://127.0.0.1:8080/target"

	http.HandleFunc("/target", func(w http.ResponseWriter, r *http.Request) {
		calledBack = true
	})

	b, err := json.Marshal(struct {
		Origin string
		Target string
	}{Origin: origin, Target: target})

	if err != nil {
		t.Fatal(err.Error())
	}

	mustHttpPost(fmt.Sprintf("http://127.0.0.1:8080/register_endpoint"), "application/json", b)

	mustHttpGet("http://127.0.0.1:8080/origin")

	if !calledBack {
		t.Error("didn't redirect to target")
	}
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

func mustHttpGet(url string) (body []byte) {
	resp, err := http.Get(url)
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

func TestGetInstanceSlices(t *testing.T) {
	var slices []string
	mock := clusterMock{iterator: iterator()}
	ts := httptest.NewServer(http.HandlerFunc(getInstanceSlices(mock)))

	body := mustHttpPost(ts.URL, "text/plain", []byte("1"))
	if err := json.Unmarshal(body, &slices); err != nil {
		t.Fatal(err.Error())
	}
	if slices[0] != "2" || slices[1] != "3" {
		t.Error(fmt.Sprintf("unexpected slices returned from server %v", slices))
	}
	defer ts.Close()
}

func TestIncrClusterId(t *testing.T) {
	mock := clusterMock{iterator: iterator()}
	mock.IncrClusterId()
	mock.IncrClusterId()

	ts := httptest.NewServer(http.HandlerFunc(incrClusterId(mock)))
	defer ts.Close()

	body := mustHttpGet(ts.URL)

	if string(body) != "2" {
		t.Error(fmt.Sprintf("unexpected id returned from server %v", string(body)))
	}

}

func startServer() {
	go main()
	runtime.Gosched()
	<-time.After(10 * time.Millisecond) //give the server some time to start
}

func (c clusterMock) GetSlices() []string {
	return []string{"0", "1", "2", "3"}
}

func iterator() func() int {
	i := -1
	return func() int {
		i += 1
		return i
	}
}
func (c clusterMock) IncrClusterId() int {
	return c.iterator()
}
func (c clusterMock) GetInstanceSlices(id int) []string {
	return c.GetSlices()[id*2:(id+1)*2]
}
