package main

import (
	"image/color"
	"math"

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

	// New constants for sword swinging
	swingDuration = 20  // frames
	swingAngleMax = 120 // degrees
)

type Game struct {
	playerX     float64
	playerY     float64
	playerVelY  float64
	onGround    bool
	platformY   float64
	facingRight bool

	// New variables for sword swinging
	swinging      bool
	swingFrame    int
	swingCooldown int

	frameLimit int
}

func NewGame() *Game {
	return &Game{
		playerX:     screenWidth / 2.0,
		playerY:     screenHeight - platformHeight - playerHeight,
		platformY:   screenHeight - platformHeight,
		onGround:    true,
		facingRight: true,
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
	if ebiten.IsKeyPressed(ebiten.KeyX) && !g.swinging && g.swingCooldown == 0 {
		g.swinging = true
		g.swingFrame = 0
	}

	// Update swing animation
	if g.swinging {
		g.swingFrame++
		if g.swingFrame >= swingDuration {
			g.swinging = false
			if g.frameLimit == 60 {
				g.swingCooldown = 15 // 30 Frames at 60fps == 0.25 seconds
			}
			//g.swingCooldown = 30 // Add a small cooldown between swings
		}
	}

	// Update cooldown
	if g.swingCooldown > 0 {
		g.swingCooldown--
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
	vector.DrawFilledRect(screen, float32(g.playerX), float32(g.playerY), playerWidth, playerHeight, color.RGBA{0, 255, 0, 255}, true)

	// Draw the sword with swing animation
	if g.swinging {
		// Calculate swing angle based on frame
		progress := float64(g.swingFrame) / float64(swingDuration)
		// Use sine function for smooth animation
		swingAngle := math.Sin(progress*math.Pi) * swingAngleMax

		// Convert angle to radians
		angleRad := swingAngle * math.Pi / 180.0

		// Calculate sword position
		var swordCenterX, swordCenterY float32
		if g.facingRight {
			// Rotate around the right side of the player
			angleRad = -angleRad // Flip angle for right-facing
			swordCenterX = float32(g.playerX + playerWidth)
			swordCenterY = float32(g.playerY + playerHeight/2)
		} else {
			// Rotate around the left side of the player
			swordCenterX = float32(g.playerX)
			swordCenterY = float32(g.playerY + playerHeight/2)
		}

		// Calculate rotated position
		cosA := float32(math.Cos(angleRad))
		sinA := float32(math.Sin(angleRad))

		// Calculate the offset to center the sword
		offsetX := -swordHeight / 2
		if !g.facingRight {
			offsetX = -swordWidth
		}

		// Draw rotated sword by rotating position around the sword's base
		x := swordCenterX + float32(offsetX)*cosA
		y := swordCenterY + float32(offsetX)*sinA
		vector.DrawFilledRect(
			screen,
			x,
			y,
			swordWidth,
			swordHeight,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}

	// Draw the platform
	vector.DrawFilledRect(screen, 0, float32(g.platformY), screenWidth, platformHeight, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic Game with Gravity and Jumping")
	ebiten.SetTPS(game.frameLimit) // Set the game loop to run at the monitor's refresh rate
	println("Monitor refresh rate:", ebiten.TPS())
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
