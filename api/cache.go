package tripapi

import (
	"time"
	"fmt"
	"io"
	"os"
	"net/http"
	"encoding/json"
	"github.com/somehibs/tripapi/util"
)

// constants
const baseUrl = "http://tripbot.tripsit.me/"
var endpoints = map[string]string {
	"alldrugs": "api/tripsit/getAllDrugs",
}
const cacheFile = "caches/%s.json"

// variables
var cache = map[string]DrugData{}

func ensureCacheDir() {
	err := os.Mkdir("caches", 0700)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		panic("Cannot ensure cache directory - are your permissions ok? ")
	}
}

type simpleError struct {
	msg string
}

func (se simpleError) Error() string {
	return fmt.Sprintf("Error during json reading: %s", se.msg)
}

func refreshCache(name, url string) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(baseUrl+url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	var out DrugReply
	fwriter := &util.FileWriter{Reader:r.Body, FileName: fmt.Sprintf(cacheFile, name)}
	defer fwriter.Close()
	out, err = decode(fwriter)
	if err != nil {
		return err
	}
	if len(out.Err) > 0 {
		return simpleError{msg: out.Err}
	}
	cache[name] = out.Data[0]
	return nil
}

func decode(reader io.Reader) (out DrugReply, err error) {
	err = json.NewDecoder(reader).Decode(&out)
	return out, err
}

func prettyTest(cacheName, drug string) {
	data := cache[cacheName]
	drugItem := data[drug]
	fmt.Printf("%+v\n\n", drugItem)
	fmt.Println(drugItem.StringProperties())
	fmt.Println("\n")
}

func refreshCaches() error {
	for name, url := range endpoints {
		err := refreshCache(name, url)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Hey, I found %d caches\r\n", len(cache))
	return nil
}

func loadFromFile(fname, name string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	out, err := decode(file)
	cache[name] = out.Data[0]
	return
}

func checkCaches() error {
	ensureCacheDir()
	for name := range endpoints {
		// Check RAM, file before triggering refresh
		if cache[name] != nil {
//			fmt.Printf("Found %s in cache\n", name)
			continue
		}
		fname := fmt.Sprintf(cacheFile, name)
		_, err := os.Stat(fname)
		if err == nil {
//			fmt.Printf("Found %s in file cache\n", name)
			err = loadFromFile(fname, name)
		}
		if err != nil {
			// Blocking cache refresh, no cache available
			fmt.Printf("Couldn't find or load %s because %s refreshing all\n", name, err)
			err = refreshCaches()
			if err != nil {
				panic(fmt.Sprintf("Failed to fetch drugs: %s", err))
			}
		}
	}
	return nil
}

func CacheTest() string {
	checkCaches()
	name := "alldrugs"
	drugItem := cache[name]["mdma"]
	fmt.Printf("Onset is a formatted field %s\n", drugItem.FormattedField("onset"))
	fmt.Printf("Marquis is a formatted field %s\n", drugItem.FormattedField("marquis"))
	prettyTest(name, "heroin")
	prettyTest(name, "dxm")
	prettyTest(name, "mdma")
	return ""
}
