package tripapi

import (
)

const drugCache = "alldrugs"

func GetDrug(name string) *Drug {
	checkCaches()
	// Check if this is actually really cheap
	dgs := cache[drugCache]
	v, simple := dgs[name]
	if simple {
		return &v
	} else {
		for _, v := range cache[drugCache] {
			for alias := range v.Aliases {
				if v.Aliases[alias] == name {
					return &v
				}
			}
		}
	}
	return nil
}
