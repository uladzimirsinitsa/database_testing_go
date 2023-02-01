package main

import (
	"github.com/joho/godotenv"
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
	"net/http"
	"bytes"
)

type Stack []string

type Data struct {
	Code int8 `json:"code"`
	Url string `json:"url"`
	Status bool `json:"status"`
	Parsing_data string `json:"parsing_data"`
}

func init() {
	godotenv.Load()
}

func (urls *Stack) IsEmpty() bool {
	return len(*urls) == 0
}

func makeRequest(url string, data []byte) bool {
	dt := bytes.NewBuffer(data)
	timeout := time.Duration(6 * time.Second)
	client := http.Client{Timeout: timeout}
	URL := os.Getenv("URL")
	response, err := client.Post(URL, "application/json", dt)
	if err != nil {
		log.Println(err)
		return false
		}
	defer response.Body.Close()
	return true
}

func (urls *Stack) Pop() (string, bool) {
	if urls.IsEmpty() {
		return "", false
	} else {
		index := len(*urls) - 1
		url := (*urls)[index]
		*urls = (*urls)[:index]
		return url, true
	}
}

func createStack() Stack {
	file, err := os.Open("urls.txt")
	if err != nil{
		log.Fatalln(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var urls Stack
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls
}

func createJsonPostRequest(url string) []byte {
	temp := &Data{
		Code: 0,
		Url: url,
		Status: false,
		Parsing_data: "parsing_data",
	}
	data, _ := json.Marshal(temp)
	return data
}

var urls = createStack()

func thread() {
	for {
		url, _ := urls.Pop()
		if url == "" {
			break
		}
		data := createJsonPostRequest(url)
		makeRequest(url, data )
	}
}

func main() {
	for i := 0; i < 160; i++	{
		go thread()
	}
	thread()
}