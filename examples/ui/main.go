package main

import (
	"log"

	"github.com/beacon/story"
)

func main() {
	game := story.NewGame(story.WithWindowSize(800, 600), story.WithWindowTitle("UI example"))

	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
