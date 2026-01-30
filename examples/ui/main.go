package main

import (
	_ "embed"
	"image"
	"image/color"
	"log"

	"github.com/beacon/story"
	"github.com/hajimehoshi/ebiten/v2"
)

func createSmoothEllipse(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 淡蓝色和白色
	lightBlue := color.RGBA{R: 173, G: 216, B: 230, A: 255}

	// 椭圆中心点和半径
	cx := float64(width) / 2
	cy := float64(height) / 2
	rx := float64(width) / 2
	ry := float64(height) / 2

	// 绘制带抗锯齿的椭圆
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算椭圆方程的值
			dx := float64(x) - cx
			dy := float64(y) - cy
			value := (dx*dx)/(rx*rx) + (dy*dy)/(ry*ry)

			// 如果完全在椭圆内
			if value <= 0.9 {
				img.Set(x, y, lightBlue)
			} else if value <= 1.1 { // 如果靠近边界，创建平滑过渡（抗锯齿）
				// 计算混合比例
				blend := 1.0 - (value-0.9)/0.2
				r := uint8(173*blend + 255*(1-blend))
				g := uint8(216*blend + 255*(1-blend))
				b := uint8(230*blend + 255*(1-blend))
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}

	return img
}

func mainUIScene() *story.Node {
	root := story.NewNode("ui_root")

	img := createSmoothEllipse(200, 80)
	btnImg := ebiten.NewImageFromImage(img)

	btnLabel := story.NewLabel("button_label", "Button1")
	btnLabel.SetPosition(0, 0)
	btnLabel.SetSize(100, 40)
	btnLabel.SetVisible(true)

	button := story.NewSprite("start_button", btnImg)
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
		story.WithRootNode(mainUIScene))

	ebiten.SetTPS(60)
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
