package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	User_id int    `json:"userId"`
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func main() {
	fmt.Println("Running getPost()")
	getPost()

	fmt.Println("\nRunning postPost()")
	postPost()
}

func postPost() {
	post := Post{
		User_id: 1,
		Title:   "Marshalling JSON",
		Body:    "Golang Marshalling Json Data!",
	}

	jsonBytes, err := json.Marshal(post)
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://jsonplaceholder.typicode.com/posts",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		log.Fatalf("request creation failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("unexpected status code: %v", resp.StatusCode)
	}

	var created Post
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&created); err != nil {
		log.Fatalf("unable to decode response: %v", err)
	}

	fmt.Println("created post id:", created.Id)
	fmt.Println("created post title:", created.Title)
	fmt.Println("created post user_id:", created.User_id)
	fmt.Println("created post body:", created.Body)
}

func getPost() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/68")
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	var post Post
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&post); err != nil {
		log.Fatalf("json decode failed: %v", err)
	}

	fmt.Println("user_id: ", post.User_id)
	fmt.Println("id: ", post.Id)
	fmt.Println("title: ", post.Title)
	fmt.Println("body: ", post.Body)
}
