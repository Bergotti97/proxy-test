// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/maemual/go-cache"

	"github.com/joho/godotenv"
)

var number int

func handler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load("cfg.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cacheSizeStr := os.Getenv("СACHE_SIZE")
	cacheSize, err := strconv.Atoi(cacheSizeStr)
	lru, err := cache.NewLRU(cacheSize)
	if err != nil {
		fmt.Println(err)
	}
	url := os.Getenv("URL")
	page, found := lru.Get(url)
	if !found {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		responseString := string(responseData)
		lru.Add(url, responseString)
		page, found = lru.Get(url)
	}
	fmt.Fprint(w, page)

}
func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	fmt.Printf("Port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
