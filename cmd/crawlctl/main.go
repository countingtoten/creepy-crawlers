package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/countingtoten/creepy-crawlers/internal/crawler"
)

func main() {
	httpClient := &http.Client{}
	c := crawler.New(httpClient)

	body, urls, err := c.Fetch("https://www.weinertworks.com")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(body)
	fmt.Println(urls)
}
