package pkg

import (
	"fmt"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var dbPath = "./data/ip2region.xdb"

func IpLocation(ip string) (location string, err error) {
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		fmt.Printf("failed to create searcher: %s\n", err.Error())
		return
	}

	defer searcher.Close()

	// do the search
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}

	return region, nil
}
