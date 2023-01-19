package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ebitick"
)

type Game struct {
	TimerSystem *ebitick.TimerSystem
	spaceTimer  *ebitick.Timer
}

func NewGame() *Game {

	game := &Game{
		TimerSystem: ebitick.NewTimerSystem(),
	}

	// Run every 60 ticks, or every second
	game.TimerSystem.AfterTicks(60, func() { fmt.Println("::") }).Loop = true

	fmt.Println("Press space to start a timer - press space again before it elapses to pause or resume it.")
	fmt.Println("Press C while the timer is running to cancel it.")

	return game

}

func (game *Game) Update() error {

	var err error

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err = errors.New("quit")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		if game.spaceTimer == nil {

			fmt.Println("Starting 3 second timer.")

			game.spaceTimer = game.TimerSystem.After(time.Second*3, func() {
				fmt.Println("The timer has elapsed.")
				game.spaceTimer = nil // Set the timer to nil so we can restart it
			})

		} else {

			if game.spaceTimer.State == ebitick.StateRunning {
				timeLeft := game.spaceTimer.TimeLeft()
				fmt.Println("The timer is now paused, with", timeLeft.ToDuration().Seconds(), "seconds /", timeLeft, "ticks left.")
				game.spaceTimer.Pause()
			} else if game.spaceTimer.State == ebitick.StatePaused {
				fmt.Println("The timer is now resumed.")
				game.spaceTimer.Resume()
			}

		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) && game.spaceTimer != nil {
		game.spaceTimer.Cancel()
		game.spaceTimer = nil // Set it to nil so we can restart it
		fmt.Println("The timer has been canceled.")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		game.TimerSystem.Clear()
		game.spaceTimer = nil
		fmt.Println("All timers canceled and removed from the TimerSystem.")
	}

	game.TimerSystem.Update()

	return err

}

func (game *Game) Draw(screen *ebiten.Image) {}

func (game *Game) Layout(w, h int) (int, int) { return 320, 240 }

func main() {

	ebiten.SetWindowTitle("Ebitick Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err.Error() != "quit" {
		panic(err)
	}

}
