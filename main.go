// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/maemual/go-cache"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cacheSize, err := strconv.Atoi(os.Getenv("СACHE_SIZE"))
	lru, err := cache.NewLRU(cacheSize)
	if err != nil {
		log.Fatal(err)
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
		//Add автоматически удалаяет при повышении лимита кэша
		lru.Add(url, responseString)
		page, found = lru.Get(url)
	}
	fmt.Fprint(w, page)
}

func main() {
	err := godotenv.Load("cfg.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
