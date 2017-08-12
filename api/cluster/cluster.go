package cluster

import (
	"fmt"
	"math/rand"
	"time"
)

type (
	Cluster interface {
		GetSlices() []string
		GetInstanceSlices(int) []string
		IncrClusterId() int
	}

	cluster struct {
		locker         *mutexLocker
		numInstances   int
		currInstanceId int
		numSlices      int
		slicer         func(int) interface{}
		randomize      func([]string) []string
	}
)

//todo: make configurable
const NUM_PARTS = 136
const NUM_INSTANCES = 1

func NewCluster() *cluster {
	rand.Seed(time.Now().UnixNano())
	return &cluster{
		locker:         newLocker(),
		currInstanceId: -1,
		numInstances:   NUM_INSTANCES,
		numSlices:      NUM_PARTS,
		slicer:         getSlice,
		randomize:      randomize,
	}
}

func (c *cluster) IncrClusterId() int {
	var res int;
	c.locker.run(func() {
		c.currInstanceId++
		res = c.currInstanceId
	})
	return res
}

func (c *cluster) GetSlices() []string {
	res := make([]string, c.numSlices)
	for i := range res {
		res[i] = fmt.Sprintf("%v", c.slicer(i))
	}
	return res
}

func (c *cluster) GetInstanceSlices(instanceId int) []string {
	slices := c.randomize(c.GetSlices())[0:c.numSlices]
	lhs := instanceId * (c.numSlices + 1) / c.numInstances
	rhs := ((instanceId + 1) * (c.numSlices + 1) / c.numInstances)
	if lhs > len(slices) {
		return []string{}
	}
	if rhs > len(slices) {
		return slices[lhs: ]
	}
	return slices[lhs: rhs]
}

func randomize(slice []string) []string {
	if len(slice) == 1 {
		return slice
	}
	for i := range slice {
		j := rand.Intn(len(slice) - 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func getSlice(pos int) interface{} {
	return fmt.Sprintf("/slices/%d.part", pos)
}
