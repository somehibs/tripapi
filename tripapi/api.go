package tripapi

import (
	"time"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

// constants
const baseUrl = "http://tripbot.tripsit.me/"
var endpoints = map[string]string {
	"alldrugs": "api/tripsit/getAllDrugs",
}

type FileWriter struct {
	reader io.Reader
	cacheName string
}

func (fw *FileWriter) Read(out []byte) (i int, e error) {
	i, e = fw.reader.Read(out)
	if i > 0 {
		ioutil.WriteFile("caches/"+fw.cacheName+".json", out[:i], 0600)
	}
	return
}

// variables
var cache = map[string]map[string]interface{}{}

func ensureCacheDir() {
	err := os.Mkdir("caches", 0600)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		panic("Cannot ensure cache directory - are your permissions ok? ")
	}
}

func refreshCaches() error {
	ensureCacheDir()
	var client = &http.Client{Timeout: 10 * time.Second}
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
