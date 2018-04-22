package tripapi

import (
	"io"
	"io/ioutil"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
)

// constants
const baseUrl = "http://tripbot.tripsit.me/"
var endpoints = map[string]string {
	"alldrugs": "api/tripsit/getAllDrugs",
}

// variables
var client = &http.Client{Timeout: 10 * time.Second}
var cache = map[string]map[string]interface{}{}

type FileWriter struct {
	reader io.Reader
	cacheName string
}

func (fw *FileWriter) Read(out []byte) (i int, e error) {
	i, e = fw.reader.Read(out)
	return
}

func refreshCaches() error {
	for name, url := range endpoints {
		r, err := client.Get(baseUrl+url)
		if err != nil {
			return err
		}
		defer r.Body.Close()
		var out map[string]interface{}
		err = json.NewDecoder(&FileWriter{reader: r.Body, cacheName: name}).Decode(&out)
		if err != nil { return err }
		cache[name] = out
	}
	fmt.Printf("Hey, I found %d caches\r\n", len(cache))
	return nil
}

func Get() string {
	err := refreshCaches()
	if err != nil {
		fmt.Printf("Failed to fetch drugs: %s", err)
	}
	return "hey"
}
