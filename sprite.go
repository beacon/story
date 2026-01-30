package story

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite 是精灵节点，嵌入Node结构
type Sprite struct {
	*Node
	image     *ebiten.Image
	mask      *image.Alpha
	drawColor color.Color
}

// NewSprite 创建新的精灵节点
func NewSprite(name string, img *ebiten.Image) *Sprite {
	b := img.Bounds()
	alphaMask := image.NewAlpha(b)
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			alphaMask.Set(i, j, img.At(i, j))
		}
	}
	sprite := &Sprite{
		Node:      NewNode(name),
		image:     img,
		mask:      alphaMask,
		drawColor: color.White,
	}

	// Clone an image but only with alpha values.
	// This is used to detect a user cursor touches the image.

	return sprite
}

func (s *Sprite) IsWithin(x, y int) bool {
	if s.mask == nil {
		return s.Node.IsWithin(x, y)
	}

	globalX, globalY := s.GetGlobalPosition()
	maskX := x - int(globalX)
	maskY := y - int(globalY)
	r, _, _, _ := s.mask.At(int(maskX), int(maskY)).RGBA()
	return r > 0
}

// SetImage 设置精灵图像
func (s *Sprite) SetImage(img *ebiten.Image) {
	s.image = img
}

// GetImage 获取精灵图像
func (s *Sprite) GetImage() *ebiten.Image {
	return s.image
}

// SetDrawColor 设置绘制颜色
func (s *Sprite) SetDrawColor(c color.Color) {
	s.drawColor = c
}

// Draw 重写绘制方法
func (s *Sprite) Draw(screen *ebiten.Image) {
	if !s.visible || s.image == nil {
		// 即使不绘制图像，也要绘制子节点
		for _, child := range s.children {
			child.Draw(screen)
		}
		return
	}

	op := &ebiten.DrawImageOptions{}

	// 计算全局位置
	gx, gy := s.GetGlobalPosition()

	// 应用变换
	// centerX := float64(s.image.Bounds().Dx()) / 2
	// centerY := float64(s.image.Bounds().Dy()) / 2
	// op.GeoM.Translate(-centerX, -centerY)
	op.GeoM.Scale(s.scaleX, s.scaleY)
	op.GeoM.Rotate(s.rotation)
	op.GeoM.Translate(gx, gy)

	// 应用颜色滤镜 - 使用 ColorScale 替代 ColorM
	if s.drawColor != color.White {
		r, g, b, a := s.drawColor.RGBA()
		op.ColorScale.Scale(
			float32(r)/65535.0,
			float32(g)/65535.0,
			float32(b)/65535.0,
			float32(a)/65535.0,
		)
	}

	screen.DrawImage(s.image, op)

	// 绘制子节点
	for _, child := range s.children {
		child.Draw(screen)
	}
}
