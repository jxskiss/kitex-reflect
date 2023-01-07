package main

import (
	"log"

	details "github.com/jxskiss/kitex-reflect/example/bookinfo/kitex_gen/cwg/bookinfo/details/detailsservice"
)

func main() {
	svr := details.NewServer(new(DetailsServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
