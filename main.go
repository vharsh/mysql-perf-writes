package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/google/gofuzz"
	"k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type PodSchema struct {
	gorm.Model
	apiversion_kind string `gorm:"type:varchar(1024)"`
	meta            string `gorm:"type:mediumtext"`
	spec            string `gorm:"type:mediumtext"`
	status          string `gorm:"type:mediumtext"`
}

func main() {
	var count, port int
	var host, user, password, dbName string
	flag.IntVar(&count, "podcount", 10, "-podcount=10")
	flag.StringVar(&user, "user", "test", "-user=maya")
	flag.StringVar(&password, "password", "test", "-password=test")
	flag.StringVar(&host, "host", "127.0.0.1", "-host=127.0.0.1")
	flag.StringVar(&dbName, "dbName", "test", "-dbName=test")
	flag.IntVar(&port, "port", 3360, "-port=3306")
	flag.Parse()
	pods := make([]v1.Pod, count)
	for i := 0; i < count; i++ {
		pods[i] = getSinglePod()
	}
	fmt.Println("Going forward to write stuff")
	persist(pods, user, password, dbName, host, port)
}

// getSinglePod returns a Pod object with fuzzed non-sense entries
func getSinglePod() v1.Pod {
	stuff := v1.Pod{}
	// Get randomness of some sort
	fuzz.NewWithSeed(time.Now().Unix())
	// TODO: Fill the structure with valid and sane items in stuff
	fuzz.New().Fuzz(&stuff)
	return stuff
}

func persist(pods []v1.Pod, user, password, dbName, host string, port int) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName))
	if err != nil {
		fmt.Printf("Error initiating a connection to DB, %v\n", err)
	}
	defer db.Close()
	db.AutoMigrate(&PodSchema{})
	for _, i := range pods {
		metaDump, err := json.Marshal(i.ObjectMeta)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		specDump, err := json.Marshal(i.Spec)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		statusDump, err := json.Marshal(i.Status)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Print("writing down stuff\n")
		db.NewRecord(&PodSchema{
			apiversion_kind: "apps/v1_pod",
			meta:            string(metaDump),
			spec:            string(specDump),
			status:          string(statusDump),
		})
		db.Create(&PodSchema{
			apiversion_kind: "apps/v1_pod",
			meta:            string(metaDump),
			spec:            string(specDump),
			status:          string(statusDump),
		})
	}
}
