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
	cache[name] = out
	return nil
}

func decode(reader io.Reader) (out DrugReply, err error) {
	err = json.NewDecoder(reader).Decode(&out)
	drug := "heroin"
	fmt.Println(out.Data[0][drug].PrettyName)
	fmt.Println(out.Data[0][drug].Onset)
	fmt.Println(out.Data[0][drug].Categories)
	fmt.Println(out.Data[0][drug].Duration)
	fmt.Println(out.Data[0][drug].Dose)
	fmt.Println(out.Data[0][drug].Aftereffects)
	return out, err
}

func refreshCaches() error {
	for name, url := range endpoints {
		refreshCache(name, url)
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
			refreshCaches()
			break
		}
	}
	return nil
}

func Get() string {
	err := checkCaches()
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch drugs: %s", err))
	}
	return ""
}
