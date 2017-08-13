package cluster

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetSlices(t *testing.T) {
	c := NewCluster()
	c.numSlices = 0
	if len(c.GetSlices()) != 0 {
		t.Error("should return empty array")
	}
	c.numSlices = 10
	slices := c.GetSlices()
	if len(slices) != 10 {
		t.Error("should return 10 items")
	}
	if slices[5] != fmt.Sprintf("/slices/5.part") {
		t.Error("returned incorrect path")
	}
}

func TestGetInstanceSlices(t *testing.T) {
	c := NewCluster()
	c.randomize = func(s []string) []string { return s }
	c.slicer = func(pos int) interface{} { return pos }

	c.numSlices = 0
	c.numInstances = 10
	for i := 0; i < c.numInstances; i++ {
		equals(t, len(c.GetInstanceSlices(i)), 0)
	}

	c.numInstances = 10
	c.numSlices = NUM_PARTS
	resStr := []string{}

	for i := 0; i < c.numInstances; i++ {
		resStr = append(resStr, c.GetInstanceSlices(i)...)
	}
	equals(t, len(resStr), c.numSlices)
	for i := 0; i < c.numSlices; i++ {
		resInt, _ := strconv.Atoi(resStr[i])
		equals(t, resInt, i)
	}

}

func equals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatal(fmt.Sprintf("expected %v to equal %v", a, b))
	}
}
