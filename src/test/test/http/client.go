package main

import (
	"fmt"
	"log"
	"net/http"
)

func req(url string, ch chan int) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("statuscode = %d\n", res.StatusCode)
	// _, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%s\n", robots)
	fmt.Printf("url [%s] loaded\n", url)
	ch <- 1
}
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go req("https://www.baidu.com", ch1)

	go req("https://www.baidu.com/?key=", ch2)

	fmt.Println("after req")

	<-ch1
	<-ch2

	fmt.Println("after req1")
}
