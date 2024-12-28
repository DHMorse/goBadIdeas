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
	minJumpSpeed   = -8.0 // Minimum jump speed when button is released quickly
)

type Game struct {
	playerX            float64
	playerY            float64
	playerVelY         float64
	onGround           bool
	platformY          float64
	facingRight        bool
	jumpingFrame       int
	isJumping          bool // New: tracks if player is currently jumping
	jumpHeld           bool // New: tracks if jump button is being held
	swinging           bool
	swingDuration      int
	swingCooldown      int
	maxJumpFrames      int
	swingCooldownFrame int
	swingFrame         int
	frameLimit         int
}

func NewGame() *Game {
	return &Game{
		playerX:            screenWidth / 2.0,
		playerY:            screenHeight - platformHeight - playerHeight,
		platformY:          screenHeight - platformHeight,
		onGround:           true,
		facingRight:        true,
		jumpingFrame:       0,
		isJumping:          false,
		jumpHeld:           false,
		swinging:           false,
		swingDuration:      ebiten.TPS() / 4,
		swingCooldown:      ebiten.TPS() / 4,
		maxJumpFrames:      ebiten.TPS() / 3,
		swingCooldownFrame: 0,
		swingFrame:         0,
		frameLimit:         ebiten.TPS(),
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

	// Jump logic
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if g.onGround {
			// Initial jump
			g.isJumping = true
			g.jumpHeld = true
			g.jumpingFrame = 0
			g.playerVelY = jumpSpeed
			g.onGround = false
		} else if g.isJumping && g.jumpHeld {
			// Continue jump while button is held
			g.jumpingFrame++
			if g.jumpingFrame < g.maxJumpFrames {
				// Gradually reduce the upward force
				g.playerVelY = jumpSpeed * float64(g.maxJumpFrames-g.jumpingFrame) /
					float64(g.maxJumpFrames)
			}
		}
	} else {
		// Jump button released
		if g.isJumping && g.jumpHeld {
			// Cut the jump short if released early
			if g.playerVelY < minJumpSpeed {
				g.playerVelY = minJumpSpeed
			}
		}
		g.jumpHeld = false
		g.isJumping = false
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
		g.isJumping = false
		g.jumpHeld = false
		g.jumpingFrame = 0
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

// Draw and Layout functions remain unchanged
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	vector.DrawFilledRect(screen,
		float32(g.playerX),
		float32(g.playerY),
		playerWidth,
		playerHeight,
		color.RGBA{0, 255, 0, 255},
		true)

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

	vector.DrawFilledRect(screen, 0, float32(g.platformY), screenWidth, platformHeight, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic Game with Variable Jump Height")
	ebiten.SetTPS(g.frameLimit)
	println("Monitor refresh rate:", ebiten.TPS(), "Hz")
	println("Swing Duration:", g.swingDuration, "frames")
	println("Swing Cooldown:", g.swingCooldown, "frames")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
