package main

import (
        "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var hostname string
var randomBytes int
var randomString string
var delay int
var listenAddr string

func generateRandomString(length int) string {
   b := make([]byte, length)
   _, err := rand.Read(b)
   if err != nil {
      panic(err)
   }
   return base64.StdEncoding.EncodeToString(b)
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(delay * int(time.Second)))
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, fmt.Sprintf("Random bytes: %d", randomBytes))
	fmt.Fprintln(w, fmt.Sprintf("Delay: %d", delay))
	fmt.Fprintln(w, randomString)
	return
}

func main() {
	host, err := os.Hostname()
	hostname = host
	if err != nil {
		panic(err)
	}
	if os.Getenv("ECHO_BYTES") != "" {
		randomBytes, _ = strconv.Atoi(os.Getenv("ECHO_BYTES"))
		randomString = generateRandomString(randomBytes)
	}
	if os.Getenv("ECHO_DELAY") != "" {
		delay, _ = strconv.Atoi(os.Getenv("ECHO_DELAY"))
	}
	listenAddr = os.Getenv("ECHO_ADDR")
	if os.Getenv("ECHO_ADDR") == "" {
		listenAddr = ":8080"
	}


	log.Printf("Listening on %s\n", listenAddr)

	http.HandleFunc("/", echoHandler)
	http.ListenAndServe(listenAddr, logMiddleware(http.DefaultServeMux))
}
