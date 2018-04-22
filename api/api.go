package tripapi

import (
	"time"
	"fmt"
	"io"
	"os"
	"net/http"
	"encoding/json"
)

// constants
const baseUrl = "http://tripbot.tripsit.me/"
var endpoints = map[string]string {
	"alldrugs": "api/tripsit/getAllDrugs",
}
const cacheFile = "caches/%s.json"

// TODO: move to another file, rename to make the passthrough more obvious
type FileWriter struct {
	reader io.Reader
	cacheName string
	file *os.File
	fileOpened bool
}

func (fw *FileWriter) Read(out []byte) (i int, e error) {
	i, e = fw.reader.Read(out)
	if i > 0 {
		if fw.fileOpened == false {
			fname := fmt.Sprintf(cacheFile, fw.cacheName)
			// This probably won't fail in an interesting way.
			os.Remove(fname)
			var err error
			fw.file, err = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			fw.fileOpened = err == nil
		}
		fw.file.Write(out[:i])
	}
	return
}

// variables
var cache = map[string]interface{}{}

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
	fwriter := &FileWriter{reader:r.Body, cacheName: name}
	defer fwriter.file.Close()
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
	data := cache[cacheName].(DrugReply).Data[0]
	drugItem := data[drug]
	fmt.Printf("%+v\n\n", drugItem)
	fmt.Println(drugItem.PrintProperties())
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
	cache[name] = out
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
			return refreshCaches()
		}
	}
	return nil
}

func Get() string {
	err := checkCaches()
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch drugs: %s", err))
	}
	name := "alldrugs"
	drugItem := cache[name].(DrugReply).Data[0]["mdma"]
	fmt.Printf("Onset is a formatted field %s\n", drugItem.FormattedField("onset"))
	fmt.Printf("Marquis is a formatted field %s\n", drugItem.FormattedField("marquis"))
	prettyTest(name, "heroin")
	prettyTest(name, "dxm")
	prettyTest(name, "mdma")
	return ""
}
