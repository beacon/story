package story

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Label struct {
	*Node
	text string
}

func NewLabel(name string, text string) *Label {
	label := &Label{
		Node: NewNode(name),
		text: text,
	}
	return label
}

func (l *Label) SetText(text string) {
	l.text = text
	l.InvalidateCache()
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) Update() {
	l.Node.Update()
}

// Draw 重写绘制方法以显示文本
func (l *Label) Draw(screen *ebiten.Image) {
	if l.text == "" {
		return
	}

	const (
		normalFontSize = 24
		bigFontSize    = 48
	)

	globalX, globalY := l.GetGlobalPosition()
	// Draw the sample text
	op := &text.DrawOptions{}
	op.GeoM.Translate(globalX, globalY)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, l.text, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op)
}
