package main

import (
	"github.com/claudiu-persoiu/go-find-subtitle/src/core"
	"github.com/claudiu-persoiu/go-find-subtitle/src/provider/opensubtitles"
	"log"
	"os"
)

func main() {

	var path string

	if len(os.Args) > 0 {
		path = os.Args[1]
	}

	if len(os.Getenv("TR_TORRENT_DIR")) > 0 && len(os.Getenv("TR_TORRENT_NAME")) > 0 {
		path = os.Getenv("TR_TORRENT_DIR") + "/" + os.Getenv("TR_TORRENT_NAME")
	}

	if len(path) == 0 {
		log.Fatalln("No file path provided")
	}

	config, err := core.GetConfig()
	if err != nil {
		panic(err)
	}

	processor := core.NewProcessor()
	openSubtitlesFinder := opensubtitles.NewFinder(config.Opensubititles, opensubtitles.NewClient(config.Opensubititles.Key))
	defer openSubtitlesFinder.Close()

	processor.AddFinder(openSubtitlesFinder)
	processor.Process(path)
}
