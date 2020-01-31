package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/google/gofuzz"
	"k8s.io/api/core/v1"
)

func main() {
	var count int
	flag.IntVar(&count, "podcount", 1000, "-podcount=1000")
	pods := make([]v1.Pod, count)
	for i := 0; i < count; i++ {
		pods[i] = getSinglePod()
	}

}

// getSinglePod returns a Pod object with fuzzed non-sense entries
func getSinglePod() v1.Pod {
	stuff := v1.Pod{}
	fuzz.NewWithSeed(time.Now().Unix())
	fuzz.New().NilChance(0).Fuzz(&stuff)
	return stuff
}
