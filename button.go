package story

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button 是按钮类，继承自ColorRectSprite
type Button struct {
	*Sprite
}

type ButtonOption func(b *Button)

func WithButtonHoverImage(img *ebiten.Image) ButtonOption {
	return func(b *Button) {
		b.mouseEventHandler.OnHover = func() {
			b.SetImage(img)
		}
	}
}

func WithButtonLeaveImage(img *ebiten.Image) ButtonOption {
	return func(b *Button) {
		b.mouseEventHandler.OnLeave = func() {
			b.SetImage(img)
		}
	}
}

func WithButtonClickImage(img *ebiten.Image) ButtonOption {
	return func(b *Button) {
		b.mouseEventHandler.OnClick = func() {
			b.SetImage(img)
		}
	}
}

// NewButton 创建新按钮
func NewButton(name string, img *ebiten.Image, mask *ebiten.Image, opts ...ButtonOption) *Button {
	button := &Button{
		Sprite: NewSprite(name, img, mask),
	}
	for _, opt := range opts {
		opt(button)
	}
	return button
}

// Update 更新按钮状态
func (b *Button) Update() {
	b.Sprite.Update()
}

// Draw 重写绘制方法以显示文本
func (b *Button) Draw(screen *ebiten.Image) {
	// 调用父类绘制方法
	b.Sprite.Draw(screen)
}
