package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/beacon/story"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed button.png
var buttonImage []byte

func mainUIScene() *story.Node {
	root := story.NewNode("ui_root")

	btnImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(buttonImage))
	if err != nil {
		log.Fatal(err)
	}

	btnLabel := story.NewLabel("button_label", "Button1")
	btnLabel.SetPosition(0, 0)
	btnLabel.SetSize(100, 40)
	btnLabel.SetVisible(true)

	button := story.NewButton("start_button", btnImg, nil)
	button.SetPosition(300, 250)
	button.SetSize(120, 60)
	button.SetVisible(true)
	button.SetMouseEventHandler(story.MouseEventHandler{
		OnClick: func() {
			btnLabel.SetText("Clicked")
		},
		OnHover: func() {
			btnLabel.SetText("Hover")
			button.SetScale(1.2, 1.2)
		},
		OnLeave: func() {
			btnLabel.SetText("Left")
			button.SetScale(1.0, 1.0)
		},
	},
	)

	button.AddChild(btnLabel)
	root.AddChild(button)

	return root
}

func main() {
	game := story.NewGame(
		story.WithWindowSize(800, 600),
		story.WithWindowTitle("UI example"),
		story.WithLogLevel("debug"),
		story.WithRootNode(mainUIScene()))

	ebiten.SetTPS(30)
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
