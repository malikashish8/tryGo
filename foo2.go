package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// PostJSON is a new post
type PostJSON struct {
	UserID int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

const baseURL string = "https://jsonplaceholder.typicode.com"

func main() {
	channel := make(chan int)
	go manipulateChannel(channel)
	a := <-channel
	fmt.Println(a)
	a = <-channel
	fmt.Println(a)
	a = <-channel
	fmt.Println(a)
	var p *PostJSON = getOnePost(1)
	fmt.Printf("Post: %+v\n", *p)

	var posts []PostJSON = getPosts()
	for _, post := range posts {
		fmt.Println(post.Id)
	}
}

func manipulateChannel(c chan int) {
	// time.Sleep(1 * time.Second)
	arr := make([]int, 10)
	for a := range arr {
		c <- a
	}
	return
}

func getOnePost(postId int) *PostJSON {
	p := PostJSON{}
	url := fmt.Sprintf(baseURL+"/posts/%v", postId)
	body := myGet(url)
	json.Unmarshal(body, &p)
	return &p
}

func getPosts() []PostJSON {
	posts := []PostJSON{}
	body := myGet(baseURL + "/posts")
	json.Unmarshal(body, &posts)
	return posts
}

func myGet(url string) []byte {
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
