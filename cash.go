package main

import (
	cache "cash/lib"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func parseEnvs() (int, int, int, string) {
	expiryEnv := os.Getenv("EXPIRY")
	cleanupEnv := os.Getenv("CLEANUP")
	portEnv := os.Getenv("PORT")
	verboseEnv := os.Getenv("VERBOSE")

	var defaultExpiry int
	var defaultCleanup int
	var port int

	if expiryEnv != "" {
		defaultExpiry = convStrToInt(expiryEnv, "Expiry")
	} else {
		defaultExpiry = 20
	}

	if cleanupEnv != "" {
		defaultCleanup = convStrToInt(cleanupEnv, "Cleanup")
	} else {
		defaultCleanup = 10
	}

	if portEnv != "" {
		port = convStrToInt(portEnv, "Port")
	} else {
		port = 9192
	}

	return defaultExpiry, defaultCleanup, port, verboseEnv
}

func convStrToInt(value, name string) int {
	output, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(name + "Must Be Int")
	}
	return output
}

func main() {

	defaultExpiry, defaultCleanup, port, verbose := parseEnvs()

	c := cache.New(defaultCleanup)

	go c.Cleanup()

	if verbose != "" {
		go c.DumpToTerminal()
	}

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
			expiry = defaultExpiry
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
			w.Write([]byte("Must supply a key in query string"))
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

	fmt.Println("💰 is up")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
