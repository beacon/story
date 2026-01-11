package story

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game 游戏主结构
type Game struct {
	root *Node
	camera    *Camera
}

type GameOption func(g *Game)

func WithWindowSize(width int, height int) GameOption {
	return func(g *Game) {
		ebiten.SetWindowSize(width, height)
	}
}

func WithWindowTitle(title string) GameOption {
	return func(g *Game) {
		ebiten.SetWindowTitle(title)
	}
}

// NewGame 创建新游戏实例
func NewGame(opts ...GameOption) *Game {
	game := &Game{
		root: NewNode("root"),
		camera:    NewCamera(),
	}

	for _, opt := range opts {
		opt(game)
	}
	
	return game
}

func (g *Game) Run() error {
	return ebiten.RunGame(g)
}

// setupExampleScene 设置示例场景
func (g *Game) setupExampleScene() {
	// 创建背景
	bg := NewColorRectSprite("Background", 800, 600, color.RGBA{50, 50, 100, 255})
	g.root.AddChild(bg)

	// 创建父容器
	container := NewNode("Container")
	container.SetPosition(200, 150)
	g.root.AddChild(container)

	// 在容器中添加几个精灵
	redRect := NewColorRectSprite("RedRect", 100, 100, color.RGBA{255, 0, 0, 255})
	redRect.SetPosition(0, 0)
	container.AddChild(redRect)

	greenRect := NewColorRectSprite("GreenRect", 80, 80, color.RGBA{0, 255, 0, 150})
	greenRect.SetPosition(20, 20)
	container.AddChild(greenRect)

	blueRect := NewColorRectSprite("BlueRect", 60, 60, color.RGBA{0, 0, 255, 200})
	blueRect.SetPosition(40, 40)
	container.AddChild(blueRect)

	// 添加独立的精灵
	yellowRect := NewColorRectSprite("YellowRect", 120, 60, color.RGBA{255, 255, 0, 200})
	yellowRect.SetPosition(400, 100)
	g.root.AddChild(yellowRect)

	// 添加旋转和缩放示例
	rotatedRect := NewColorRectSprite("RotatedRect", 80, 40, color.RGBA{255, 0, 255, 200})
	rotatedRect.SetPosition(300, 300)
	rotatedRect.SetRotation(0.5) // 弧度
	rotatedRect.SetScale(1.5, 0.7)
	g.root.AddChild(rotatedRect)
}

// Update 更新游戏逻辑
func (g *Game) Update() error {
	g.root.Update()
	return nil
}

// Draw 绘制游戏画面
func (g *Game) Draw(screen *ebiten.Image) {
	// 清空屏幕
	screen.Fill(color.RGBA{30, 30, 50, 255})

	// 绘制场景中的所有节点
	g.root.Draw(screen)

	// 显示调试信息
	ebitenutil.DebugPrint(screen, "Godot-style Node System with Ebiten\nUse nodes as base class for all game objects")
}

// Layout 设置游戏布局
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
