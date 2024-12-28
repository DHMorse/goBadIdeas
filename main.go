package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth    = 800
	screenHeight   = 600
	playerWidth    = 50
	playerHeight   = 80
	platformHeight = 20
	gravity        = 1.0
	jumpSpeed      = -15.0
	moveSpeed      = 5.0
	swordWidth     = 80
	swordHeight    = 20
)

type Game struct {
	playerX     float64
	playerY     float64
	playerVelY  float64
	onGround    bool
	platformY   float64
	facingRight bool

	// New variables for sword swinging
	swinging           bool
	swingDuration      int
	swingCooldown      int
	swingCooldownFrame int
	swingFrame         int

	frameLimit int
}

func NewGame() *Game {
	return &Game{
		playerX:     screenWidth / 2.0,
		playerY:     screenHeight - platformHeight - playerHeight,
		platformY:   screenHeight - platformHeight,
		onGround:    true,
		facingRight: true,

		// New variables for sword swinging
		swinging:           false,
		swingDuration:      ebiten.TPS() / 4, // 0.25 seconds
		swingCooldown:      ebiten.TPS() / 4, // 0.25 seconds
		swingCooldownFrame: 0,
		swingFrame:         0,

		// get the users montior refresh rate
		frameLimit: ebiten.TPS(),
	}
}

func (g *Game) Update() error {
	// Horizontal movement
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += moveSpeed
		g.facingRight = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= moveSpeed
		g.facingRight = false
	}

	// Jump
	if g.onGround && ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.playerVelY = jumpSpeed
		g.onGround = false
	}

	// Sword swing
	if ebiten.IsKeyPressed(ebiten.KeyX) && !g.swinging && g.swingCooldownFrame == 0 {
		g.swinging = true
		g.swingFrame = 0
	}

	// Update swing animation
	if g.swinging {
		g.swingFrame++
		if g.swingFrame >= g.swingDuration {
			g.swinging = false
			g.swingCooldownFrame = g.swingDuration
		}
	}

	// Update cooldown
	if g.swingCooldownFrame > 0 {
		g.swingCooldownFrame--
	}

	// Gravity
	g.playerVelY += gravity
	g.playerY += g.playerVelY

	// Collision with the ground/platform
	if g.playerY+playerHeight >= g.platformY {
		g.playerY = g.platformY - playerHeight
		g.playerVelY = 0
		g.onGround = true
	}

	// Keep the player within the screen bounds
	if g.playerX < 0 {
		g.playerX = 0
	}
	if g.playerX+playerWidth > screenWidth {
		g.playerX = screenWidth - playerWidth
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with a color
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw the player
	vector.DrawFilledRect(screen,
		float32(g.playerX),
		float32(g.playerY),
		playerWidth,
		playerHeight,
		color.RGBA{0, 255, 0, 255},
		true)

	// Draw the sword with swing animation
	if g.swinging {
		if g.facingRight {
			vector.DrawFilledRect(
				screen,
				float32(g.playerX)+playerWidth,
				float32(g.playerY)+playerHeight/2-swordHeight/2,
				swordWidth,
				swordHeight,
				color.RGBA{255, 0, 0, 255},
				true,
			)
		} else {
			vector.DrawFilledRect(
				screen,
				float32(g.playerX)-swordWidth,
				float32(g.playerY)+playerHeight/2-swordHeight/2,
				swordWidth,
				swordHeight,
				color.RGBA{255, 0, 0, 255},
				true,
			)
		}
	}

	// Draw the platform
	vector.DrawFilledRect(screen, 0, float32(g.platformY), screenWidth, platformHeight, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic Game with Gravity and Jumping")
	ebiten.SetTPS(g.frameLimit) // Set the game loop to run at the monitor's refresh rate
	println("Monitor refresh rate:", ebiten.TPS(), "Hz")
	println("Swing Duration:", g.swingDuration, "frames")
	println("Swing Cooldown:", g.swingCooldown, "frames")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
