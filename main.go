package main

import (
	"fmt"
	"github.com/somehibs/tripapi/api"
)

func main() {
	fmt.Println(tripapi.CacheTest()) // cache test
	fmt.Println(tripapi.GetDrug("heroin")) // heroin itself
	fmt.Println(tripapi.GetDrug("dex")) // dxm alias
	fmt.Println(tripapi.GetDrug("butt")) // nonexistent
}
