# ebitick⏱️

[![Go Reference](https://pkg.go.dev/badge/github.com/solarlune/ebitick.svg)](https://pkg.go.dev/github.com/solarlune/ebitick)

ebitick is a timer system for [Ebitengine](https://ebitengine.org/) games.

## Why?

Because timing stuff is important and can be done in an easy to use, set-it-and-forget-it kinda way. You can just use `time.After()`, but that works on an additional goroutine, which can introduce race conditions into the mix. Ebitick, in comparison, works on the same goroutine as the rest of your game.

❗ : Note that ebitick keeps time by counting ticks against the target tickrate (TPS), so it won't work properly if you change the tickrate while timers are running.

## How?

`go get github.com/solarlune/ebitick`

```go

package main

import (
    "fmt"
    "time"
    
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/ebitick"
)

type Game struct {
    TimerSystem *ebitick.TimerSystem
}

func NewGame() *Game {

    game := &Game{
        TimerSystem: ebitick.NewTimerSystem(),
    }

    // Below, we specify that after a second has passed,
    // we run the given function literal.
    game.TimerSystem.After(time.Second, func() {
        fmt.Println("A second has passed.")
    })

    // There's also a TimerSystem.AfterTicks() function if you want to use
    // ticks exactly without any human-readable time.Duration conversions.
    
    // The After functions return the created Timer, allowing you to pause it
    // later as necessary or set it to loop, for example. Once the timer elapses,
    // it is removed from the TimerSystem and its state is set to StateFinished,
    // so you can check for that if necessary as well.

    return game

}

func (game *Game) Update() error {

    // We have the TimerSystem update once per game tick, and that's it.
    game.TimerSystem.Update()

    return nil

}

func (game *Game) Draw(screen *ebiten.Image) {}

func (game *Game) Layout(w, h int) (int, int) { return 320, 240 }

func main() {

    ebiten.SetWindowTitle("Ebitick Minimal Example")
    ebiten.RunGame(NewGame())

}

```

## Shout-outs

[Ebitengine's Discord server~](https://discord.gg/fXM7VYASTu)
