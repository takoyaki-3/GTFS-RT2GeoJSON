package main

import (
	"fmt"
	"github.com/takoyaki-3/GTFS-RT2GeoJSON/pkg"
)

func main(){
	fmt.Println("start")

	paths,_ := pkg.Dirwalk("./zip")

	for _,v:=range paths{
		fmt.Println(v)

		pkg.Unzip(v,"./GTFS-RTs")
	}
}

