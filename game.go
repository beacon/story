package story

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sirupsen/logrus"
)

// Game 游戏主结构
type Game struct {
	root   *Node
	camera *Camera

	width           int
	height          int
	backgroundColor color.Color
}

type GameOption func(g *Game)

func WithLogLevel(lvl string) GameOption {
	return func(g *Game) {
		level, err := logrus.ParseLevel(lvl)
		if err != nil {
			logrus.Error("Invalid log level:", err)
			return
		}
		logrus.SetLevel(level)
	}
}

func WithWindowSize(width int, height int) GameOption {
	return func(g *Game) {
		g.width = width
		g.height = height
		ebiten.SetWindowSize(width, height)
	}
}

func WithWindowTitle(title string) GameOption {
	return func(g *Game) {
		ebiten.SetWindowTitle(title)
	}
}

func WithBackgroundColor(col color.Color) GameOption {
	return func(g *Game) {
		g.backgroundColor = col
	}
}

func WithRootNode(root *Node) GameOption {
	return func(g *Game) {
		g.root = root
	}
}

// NewGame 创建新游戏实例
func NewGame(opts ...GameOption) *Game {
	game := &Game{
		root:            NewNode("root"),
		camera:          NewCamera(),
		backgroundColor: color.RGBA{},
	}

	for _, opt := range opts {
		opt(game)
	}

	return game
}

func (g *Game) Run() error {
	return ebiten.RunGame(g)
}

// Update 更新游戏逻辑
func (g *Game) Update() error {
	g.root.Update()
	return nil
}

// Draw 绘制游戏画面
func (g *Game) Draw(screen *ebiten.Image) {
	// 清空屏幕
	screen.Fill(g.backgroundColor)

	// 绘制场景中的所有节点
	g.root.Draw(screen)

	// 显示调试信息
	ebitenutil.DebugPrint(screen, "Godot-style Node System with Ebiten\nUse nodes as base class for all game objects")
}

// Layout 设置游戏布局
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
