package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// Button 是按钮类，继承自ColorRectSprite
type Button struct {
	*Sprite
	text          string
	fontFace      font.Face
	onClick       EventCallback
	onHoverEnter  EventCallback
	onHoverLeave  EventCallback
	isHovered     bool
	defaultColor  color.Color
	hoverColor    color.Color
	clickColor    color.Color
	pressed       bool
}

// NewButton 创建新按钮
func NewButton(name string, image string) *Button {
	button := &Button{
		Sprite: NewSprite(name, image),
		text:            text,
		fontFace:        basicfont.Face7x13,
		defaultColor:    defaultColor,
		hoverColor:      hoverColor,
		clickColor:      clickColor,
		isHovered:       false,
		pressed:         false,
	}
	
	return button
}

// SetText 设置按钮文本
func (b *Button) SetText(text string) {
	b.text = text
	b.InvalidateCache()
}

// OnClick 注册点击事件处理器
func (b *Button) OnClick(callback EventCallback) {
	b.onClick = callback
}

// OnHoverEnter 注册鼠标进入事件处理器
func (b *Button) OnHoverEnter(callback EventCallback) {
	b.onHoverEnter = callback
}

// OnHoverLeave 注册鼠标离开事件处理器
func (b *Button) OnHoverLeave(callback EventCallback) {
	b.onHoverLeave = callback
}

// Update 更新按钮状态
func (b *Button) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	globalX, globalY := b.GetGlobalPosition()
	
	// 检查鼠标是否悬停在按钮上
	isOver := mouseX >= int(globalX) && mouseX < int(globalX)+b.width &&
	          mouseY >= int(globalY) && mouseY < int(globalY)+b.height
	
	// 处理悬停状态
	if isOver && !b.isHovered {
		// 鼠标进入
		b.isHovered = true
		if b.onHoverEnter != nil {
			b.onHoverEnter(nil)
		}
		b.ChangeColor(b.hoverColor)
	} else if !isOver && b.isHovered {
		// 鼠标离开
		b.isHovered = false
		if b.onHoverLeave != nil {
			b.onHoverLeave(nil)
		}
		b.ChangeColor(b.defaultColor)
		b.pressed = false
	}
	
	// 检查点击
	if isOver && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.pressed = true
		b.ChangeColor(b.clickColor)
		if b.onClick != nil {
			b.onClick(nil)
		}
	}
	
	// 如果按钮被按下但鼠标移出了按钮区域，恢复到悬停状态
	if b.pressed && !isOver && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.ChangeColor(b.hoverColor)
	} else if b.pressed && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// 如果按钮被按下但松开了鼠标，恢复到悬停或默认状态
		b.pressed = false
		if isOver {
			b.ChangeColor(b.hoverColor)
		} else {
			b.ChangeColor(b.defaultColor)
		}
	}
	
	// 更新子节点
	for _, child := range b.GetChildren() {
		child.Update()
	}
}

// ChangeColor 改变按钮颜色
func (b *Button) ChangeColor(newColor color.Color) {
	b.UpdateRect(newColor)
}

// Draw 重写绘制方法以显示文本
func (b *Button) Draw(screen *ebiten.Image) {
	// 调用父类绘制方法
	b.ColorRectSprite.Draw(screen)
	
	// 绘制文本
	if b.text != "" {
		globalX, globalY := b.GetGlobalPosition()
		textX := globalX + float64(b.width)/2
		textY := globalY + float64(b.height)/2 + 5 // 垂直居中调整
		
		// 计算文本宽度以水平居中文本
		textWidth := float64(len(b.text)) * 7 // 基本字体的近似宽度
		textX -= textWidth / 2
		
		ebitenutil.DebugPrintAt(screen, b.text, int(textX), int(textY))
	}
}
