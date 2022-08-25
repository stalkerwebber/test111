package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var Collection map[string]int
var CollectionMux sync.RWMutex

func routes() {
	router.HandleFunc("/", Home)
	router.HandleFunc("/ping", Ping)
	router.HandleFunc("/collect", Collect)
	router.HandleFunc("/stat", Stat)

	Collection = make(map[string]int)
}

func Home(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", "usage: http://127.0.0.1:4100/collect?words=qwe,asd http://127.0.0.1:4100/stat")
}

func Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "pong")
}

func Collect(w http.ResponseWriter, r *http.Request) {
	words := urlParam(r, "words")
	wordsArr := strings.Split(words, ",")
	CollectionMux.Lock()
	defer CollectionMux.Unlock()
	for _, word := range wordsArr {
		if len(word) > 0 {
			Collection[word] += 1
		}
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("collected: %s", wordsArr)
	log.Printf("collection: %+v", Collection)
}

func Stat(w http.ResponseWriter, r *http.Request) {
	// get distribution rank
	var CollectionCopy map[string]int
	CollectionCopy = make(map[string]int)
	// make copy to release collection as soon as possible
	CollectionMux.RLock()
	var min, max int
	for word, qty := range Collection {
		if qty < min {
			min = qty
		}
		if qty > max {
			max = qty
		}
		CollectionCopy[word] = qty
	}
	CollectionMux.RUnlock()

	valuesRange := max - min

	for word, qty := range CollectionCopy {
		CollectionCopy[word] = int(math.Round(float64(4*(qty-min)/valuesRange)) + 1)
	}

	keys := make([]string, 0, len(CollectionCopy))

	for key := range CollectionCopy {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return CollectionCopy[keys[i]] > CollectionCopy[keys[j]]
	})

	log.Printf("%+x", CollectionCopy)

	w.WriteHeader(http.StatusOK)
	i := 5
	for _, word := range keys {
		if i <= 0 {
			break
		}
		_, _ = fmt.Fprintf(w, "%s %d\n", word, CollectionCopy[word])
		i--
	}
}

func urlParam(r *http.Request, paramName string) string {
	vs, ok := r.URL.Query()[paramName]
	if ok {
		return vs[0]
	}
	return ""
}

func responseError(err string, w http.ResponseWriter) {
	w.Header().Add("Error", err)
	w.Write([]byte("Error: " + err))
	w.WriteHeader(http.StatusBadRequest)
}
