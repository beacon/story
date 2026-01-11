package story

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type EventCallback func(interface{})

// RenderCache 缓存渲染结果
type RenderCache struct {
	image       *ebiten.Image
	dirty       bool
	lastBounds  [4]int // x, y, width, height
}

// NewRenderCache 创建新的渲染缓存
func NewRenderCache(width, height int) *RenderCache {
	return &RenderCache{
		image: ebiten.NewImage(width, height),
		dirty: true,
	}
}

// NodeInterface 定义节点的基本接口
type NodeInterface interface {
	GetName() string
	SetName(string)
	GetVisible() bool
	SetVisible(bool)
	GetPosition() (float64, float64)
	SetPosition(float64, float64)
	GetRotation() float64
	SetRotation(float64)
	GetScale() (float64, float64)
	SetScale(float64, float64)
	GetParent() NodeInterface
	SetParent(NodeInterface)
	GetChildren() []NodeInterface
	AddChild(NodeInterface)
	RemoveChild(NodeInterface)
	GetGlobalPosition() (float64, float64)
	Update()
	Draw(*ebiten.Image)
	IsStatic() bool
	SetStatic(bool)
	InvalidateCache()
	GetCache() *RenderCache
	SetCache(*RenderCache)
}

// Node 是基础节点实现
type Node struct {
	name        string
	visible     bool
	x, y        float64
	rotation    float64 // 旋转角度（弧度）
	scaleX      float64 // X轴缩放
	scaleY      float64 // Y轴缩放
	parent      NodeInterface
	children    []NodeInterface
	static      bool           // 是否为静态节点
	cache       *RenderCache   // 渲染缓存
	mutex       sync.RWMutex   // 保护并发访问
}

// NewNode 创建新的基础节点
func NewNode(name string) *Node {
	return &Node{
		name:    name,
		visible: true,
		scaleX:  1.0,
		scaleY:  1.0,
		static:  false,
	}
}

// 接口实现
func (n *Node) GetName() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.name
}

func (n *Node) SetName(name string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.name = name
}

func (n *Node) GetVisible() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.visible
}

func (n *Node) SetVisible(visible bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.visible = visible
	n.InvalidateCache()
}

func (n *Node) GetPosition() (float64, float64) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.x, n.y
}

func (n *Node) SetPosition(x, y float64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.x = x
	n.y = y
	n.InvalidateCache()
}

func (n *Node) GetRotation() float64 {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.rotation
}

func (n *Node) SetRotation(rotation float64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	oldRotation := n.rotation
	n.rotation = rotation
	if oldRotation != rotation {
		n.InvalidateCache()
	}
}

func (n *Node) GetScale() (float64, float64) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.scaleX, n.scaleY
}

func (n *Node) SetScale(scaleX, scaleY float64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	oldScaleX, oldScaleY := n.scaleX, n.scaleY
	n.scaleX = scaleX
	n.scaleY = scaleY
	if oldScaleX != scaleX || oldScaleY != scaleY {
		n.InvalidateCache()
	}
}

func (n *Node) GetParent() NodeInterface {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.parent
}

func (n *Node) SetParent(parent NodeInterface) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.parent = parent
}

func (n *Node) GetChildren() []NodeInterface {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.children
}

func (n *Node) AddChild(child NodeInterface) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	child.SetParent(n)
	n.children = append(n.children, child)
	n.InvalidateCache()
}

func (n *Node) RemoveChild(child NodeInterface) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	for i, c := range n.children {
		if c == child {
			c.SetParent(nil)
			n.children = append(n.children[:i], n.children[i+1:]...)
			n.InvalidateCache()
			break
		}
	}
}

func (n *Node) GetGlobalPosition() (float64, float64) {
	n.mutex.RLock()
	x, y := n.x, n.y
	parent := n.parent
	n.mutex.RUnlock()
	
	for parent != nil {
		px, py := parent.GetPosition()
		x += px
		y += py
		parent = parent.GetParent()
	}
	return x, y
}

func (n *Node) IsStatic() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.static
}

func (n *Node) SetStatic(static bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	oldStatic := n.static
	n.static = static
	if oldStatic != static {
		n.InvalidateCache()
	}
}

func (n *Node) InvalidateCache() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.cache != nil {
		n.cache.dirty = true
	}
	// 同时使所有子节点的缓存失效
	for _, child := range n.children {
		child.InvalidateCache()
	}
}

func (n *Node) GetCache() *RenderCache {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.cache
}

func (n *Node) SetCache(cache *RenderCache) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.cache = cache
}

func (n *Node) Update() {
	n.mutex.RLock()
	isStatic := n.static
	n.mutex.RUnlock()
	
	// 如果是静态节点且有缓存，则不需要频繁更新
	if !isStatic {
		for _, child := range n.children {
			child.Update()
		}
	} else {
		// 对于静态节点，我们仍然需要更新子节点，但如果子节点不是静态的则更新
		for _, child := range n.children {
			if !child.IsStatic() {
				child.Update()
			}
		}
	}
}

func (n *Node) Draw(screen *ebiten.Image) {
	n.mutex.RLock()
	visible := n.visible
	static := n.static
	n.mutex.RUnlock()
	
	if !visible {
		return
	}
	
	// 尝试使用缓存
	var renderTarget *ebiten.Image
	
	if static && n.cache != nil && !n.cache.dirty {
		// 使用缓存
		renderTarget = n.cache.image
	} else if static {
		// 需要重新渲染到缓存
		n.renderToCache()
		renderTarget = n.cache.image
	} else {
		// 动态节点直接绘制
		renderTarget = nil
	}
	
	if renderTarget != nil {
		// 使用缓存绘制
		op := &ebiten.DrawImageOptions{}
		
		gx, gy := n.GetGlobalPosition()
		
		// 应用变换
		op.GeoM.Scale(n.scaleX, n.scaleY)
		op.GeoM.Rotate(n.rotation)
		op.GeoM.Translate(gx, gy)
		
		screen.DrawImage(renderTarget, op)
	} else {
		// 动态绘制 - 先绘制自己（如果需要），然后绘制子节点
		n.drawChildren(screen)
	}
}

func (n *Node) renderToCache() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	
	if n.cache == nil {
		// 估算缓存大小，这里使用固定大小作为示例
		n.cache = NewRenderCache(800, 600)
	}
	
	// 清空缓存
	n.cache.image.Fill(color.Transparent)
	
	// 渲染子节点到缓存
	n.drawChildren(n.cache.image)
	
	n.cache.dirty = false
}

func (n *Node) drawChildren(screen *ebiten.Image) {
	n.mutex.RLock()
	children := make([]NodeInterface, len(n.children))
	copy(children, n.children)
	visible := n.visible
	n.mutex.RUnlock()
	
	if !visible {
		return
	}
	
	// 绘制所有子节点
	for _, child := range children {
		child.Draw(screen)
	}
}

// Sprite 是精灵节点，嵌入Node结构
type Sprite struct {
	*Node
	image     *ebiten.Image
	drawColor color.Color
}

// NewSprite 创建新的精灵节点
func NewSprite(name string, img *ebiten.Image) *Sprite {
	sprite := &Sprite{
		Node:      NewNode(name),
		image:     img,
		drawColor: color.White,
	}
	return sprite
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
	centerX := float64(s.image.Bounds().Dx()) / 2
	centerY := float64(s.image.Bounds().Dy()) / 2
	op.GeoM.Translate(-centerX, -centerY)
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

// ColorRectSprite 是彩色矩形精灵
type ColorRectSprite struct {
	*Sprite
	width, height int
	rectColor     color.Color
}

// NewColorRectSprite 创建彩色矩形精灵
func NewColorRectSprite(name string, width, height int, rectColor color.Color) *ColorRectSprite {
	img := ebiten.NewImage(width, height)
	img.Fill(rectColor)
	
	sprite := &ColorRectSprite{
		Sprite:    NewSprite(name, img),
		width:     width,
		height:    height,
		rectColor: rectColor,
	}
	
	return sprite
}

// UpdateRect 更新矩形颜色
func (crs *ColorRectSprite) UpdateRect(newColor color.Color) {
	crs.rectColor = newColor
	crs.image.Fill(crs.rectColor)
}

// UpdateSize 更新矩形大小
func (crs *ColorRectSprite) UpdateSize(width, height int) {
	crs.width = width
	crs.height = height
	newImg := ebiten.NewImage(width, height)
	newImg.Fill(crs.rectColor)
	crs.image = newImg
}

// Camera 相机系统（简化版）
type Camera struct {
	x, y         float64
	zoom         float64
	rotation     float64
	followTarget NodeInterface
}

// NewCamera 创建新相机
func NewCamera() *Camera {
	return &Camera{
		zoom: 1.0,
	}
}

// GetPosition 获取相机位置
func (c *Camera) GetPosition() (float64, float64) {
	return c.x, c.y
}

// SetPosition 设置相机位置
func (c *Camera) SetPosition(x, y float64) {
	c.x = x
	c.y = y
}

// GetZoom 获取相机缩放
func (c *Camera) GetZoom() float64 {
	return c.zoom
}

// SetZoom 设置相机缩放
func (c *Camera) SetZoom(zoom float64) {
	c.zoom = zoom
}

// Follow 设置相机跟随目标
func (c *Camera) Follow(target NodeInterface) {
	c.followTarget = target
}

// Update 更新相机逻辑
func (c *Camera) Update() {
	if c.followTarget != nil {
		x, y := c.followTarget.GetGlobalPosition()
		c.x = x
		c.y = y
	}
}

// ApplyToOptions 应用相机变换到绘图选项
func (c *Camera) ApplyToOptions(op *ebiten.DrawImageOptions) {
	// 实现相机变换逻辑
	op.GeoM.Translate(-c.x, -c.y)
	op.GeoM.Scale(c.zoom, c.zoom)
	op.GeoM.Rotate(c.rotation)
}
