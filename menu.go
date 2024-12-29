package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	X, Y          float64
	Width, Height float64
	Text          string
	OnClick       func()
}

func (b *Button) Draw(screen *ebiten.Image, isHovered bool) {
	// Change button color when hovered
	buttonColor := color.RGBA{100, 100, 100, 255}
	if isHovered {
		buttonColor = color.RGBA{150, 150, 150, 255}
	}

	// Draw button background
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), buttonColor, true)

	// Draw button text
	textX := int(b.X + b.Width/4) // Adjust text position to center
	textY := int(b.Y + b.Height/4)
	ebitenutil.DebugPrintAt(screen, b.Text, textX, textY)
}

func (b *Button) IsMouseOver(x, y float64) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

func DrawMainMenu(screen *ebiten.Image, g *Game) {
	buttons := []*Button{
		{
			X:      100,
			Y:      100,
			Width:  200,
			Height: 50,
			Text:   "Play",
			OnClick: func() {
				// Transition to the game
				g.state.MainMenu = false
				g.state.Playing = true
			},
		},
		{
			X:      100,
			Y:      200,
			Width:  200,
			Height: 50,
			Text:   "Settings",
			OnClick: func() {
				// Add settings functionality
				println("Settings Clicked")
			},
		},
		{
			X:      100,
			Y:      300,
			Width:  200,
			Height: 50,
			Text:   "Quit",
			OnClick: func() {
				// Exit the game
				println("Quit Clicked")
			},
		},
	}

	// Get mouse position
	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float64(mx), float64(my)

	// Draw buttons and check for clicks
	for _, button := range buttons {
		isHovered := button.IsMouseOver(mouseX, mouseY)
		button.Draw(screen, isHovered)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && isHovered {
			button.OnClick()
		}
	}
}
