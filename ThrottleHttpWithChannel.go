package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// PostJSON is a new post
type PostJSON struct {
	UserID int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

const baseURL string = "http://jsonplaceholder.typicode.com"

const RequestsPerSecond int = 13

var throttle chan int = make(chan int, RequestsPerSecond)

func main() {
	go operateThrottle()

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8080")
	fmt.Println("Using Proxy")
	channel := make(chan PostJSON)
	var posts []PostJSON = getPosts()
	for _, post := range posts {
		go getOnePost(post.Id, channel)
	}

	for i := range posts {
		var postDetails PostJSON = <-channel
		fmt.Println(i, postDetails.Id)
	}
}

func operateThrottle() {
	iterator := make([]int, cap(throttle))
	for true {
		time.Sleep(1 * time.Second)
		fmt.Printf("Processed %v in last second\n", len(throttle))
		for range iterator[0:len(throttle)] {
			<-throttle
		}
	}
}

func getOnePost(postId int, channel chan PostJSON) {
	p := PostJSON{}
	url := fmt.Sprintf(baseURL+"/posts/%v", postId)
	body := myGet(url)
	err := json.Unmarshal(body, &p)
	if err != nil {
		handleError(err)
	}
	channel <- p
	return
}

func getPosts() []PostJSON {
	posts := []PostJSON{}
	body := myGet(baseURL + "/posts")
	json.Unmarshal(body, &posts)
	return posts
}

func myGet(url string) []byte {
	throttle <- 1
	fmt.Println("GETting URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		handleError(err)
		return nil
	} else {
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body)
		// fmt.Printf("Read Body:\n%+v\n", string(body))
		if err2 != nil {
			handleError(err)
		}
		return body
	}
}

func handleError(err error) {
	fmt.Println("Error occured\n", err)
	os.Exit(1)
}
