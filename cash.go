package main

import (
	cache "cash/lib"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	defaultExpiry := flag.Int("expiry", 20, "Default expiration time of items in cache.")
	defaultCleanup := flag.Int("cleanup", 10, "Time between cleanups.")
	defaultPort := flag.Int("port", 9192, "The port the service runs on.")
	// verbose := flag.Bool("verbose", false, "Verbose mode to output cache to terminal every 2 minutes.")

	flag.Parse()

	// defaultExpiry := os.Getenv("EXPIRY")
	// defaultCleanup := os.Getenv("CLEANUP")
	// defaultPort := os.Getenv("PORT")
	// verbose := os.Getenv("VERBOSE")

	// if defaultCleanup != "" {
	// 	log.Fatal("EXPIRY must be an to an int")
	// }

	c := cache.New(*defaultCleanup)

	go c.Cleanup()
	// go c.DumpToTerminal()

	// if verbose != "" {
	// 	go c.DumpToTerminal()
	// }

	// Handle Create
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		value := r.URL.Query().Get("value")

		if value == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Must supply a value in query string"))
			return
		}

		expiry, err := strconv.Atoi(r.URL.Query().Get("expiry"))

		if err != nil || expiry == 0 {
			expiry = *defaultExpiry
		}

		hash := md5.New()
		hash.Write([]byte(value + startTime.String()))
		hashedKey := hex.EncodeToString(hash.Sum(nil))

		c.CreateItem(hashedKey, value, expiry)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hashedKey))
	})

	// Handle Read
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Query().Get("key")

		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Must supply a value in query string"))
			return
		}

		value, err := c.ReadItem(key)

		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("Not Found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	})

	// Handle Destroy
	http.HandleFunc("/destroy", func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Query().Get("key")

		c.DestroyItem(key)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Destroyed"))
	})

	// Handle Purge
	http.HandleFunc("/purge", func(w http.ResponseWriter, r *http.Request) {
		c.Purge()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Purged"))
	})

	// Handle Health Check
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("CASH IS UP")
	http.ListenAndServe(":"+strconv.Itoa(*defaultPort), nil)
}
