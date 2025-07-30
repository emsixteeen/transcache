package main

import (
	"fmt"
	"github.com/emsixteeen/transcache"
	"github.com/jessevdk/go-flags"
	"log"
)

var args struct {
	FFMpeg     string `short:"m" long:"ffmpeg" default:"ffmpeg"`
	ListenAddr string `short:"l" long:"listen" default:":9999" required:"true"`
}

func main() {
	if _, err := flags.Parse(&args); err != nil {
		return
	}

	fmt.Println("ffmpeg:", args.FFMpeg)
	fmt.Println("listen:", args.ListenAddr)

	srvr := transcache.Server{
		Addr: args.ListenAddr,
		Converter: transcache.Converter{
			Exec: args.FFMpeg,
		},
	}

	log.Fatal(srvr.Run())
}
