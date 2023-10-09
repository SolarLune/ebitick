package ebitick

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// TimeUnit represents a game tick in an ebitengine game. For simplicity, a TimeUnit can be used as either a timestamp
// (think time.Time{}, time.Now()), or a duration of time (time.Duration{}, time.Since()) depending on the context with
// which the value is used. It is a float so that a TimerSystem can run at faster or slower speeds.
type TimeUnit float32

// ToDuration converts the timestamp to a generic time.Duration.
func (ts TimeUnit) ToDuration() time.Duration {
	return time.Duration(float64(ts) / float64(ebiten.TPS()) * float64(time.Second))
}

// ToTimeUnit converts the given number of seconds to a TimeUnit using Ebiten's current TPS value.
func ToTimeUnit(duration time.Duration) TimeUnit {
	return TimeUnit(duration.Seconds() * float64(ebiten.TPS()))
}

// The various possible states for a Timer.
const (
	StateRunning = iota
	StateCanceled
	StatePaused
	StateFinished
)

// Timer represents a timer that elapses after a given amount of time.
type Timer struct {
	timerSystem *TimerSystem
	StartTick   TimeUnit // On what tick of the TimerSystem the Timer was initially started.
	duration    TimeUnit // How long the Timer should take.
	OnExecute   func()   // What the timer does once it elapses.
	Loop        bool     // If the Timer should loop after elapsing. Defaults to off.
	State       int      // What state the Timer is in.
}

// Cancel cancels a Timer, removing it from the TimerSystem the next time TimerSystem.Update() is called. This does nothing on a finished Timer.
func (timer *Timer) Cancel() {
	if timer.State != StateFinished {
		timer.State = StateCanceled
		timer.timerSystem.removeTimer(timer)
	}
}

// Pause pauses the Timer. While paused, a Timer is not incrementing time. This does nothing on a Timer if it isn't running, specifically.
func (timer *Timer) Pause() {
	if timer.State == StateRunning {
		timer.State = StatePaused
	}
}

// Resume resumes a paused Timer. This does nothing on a Timer if it isn't paused, specifically.
func (timer *Timer) Resume() {
	if timer.State == StatePaused {
		timer.State = StateRunning
	}
}

// TimeLeft returns a TimeUnit indicating how much -absolute- time is left on the Timer. This value is multiplied
// by the owning system's current speed value.
func (timer *Timer) TimeLeft() TimeUnit {
	return ((timer.duration + timer.StartTick) - timer.timerSystem.CurrentTime) / TimeUnit(timer.timerSystem.Speed)
}

func (timer *Timer) SetDuration(duration TimeUnit) {
	timer.duration = duration
}

func (timer *Timer) Restart() {
	timer.StartTick = timer.timerSystem.CurrentTime
}

// TimerSystem represents a system that updates and triggers timers added to the System.
type TimerSystem struct {
	Timers      []*Timer // The Timers presently existing in the System.
	CurrentTime TimeUnit // The current TimeUnit (tick) of the TimerSystem. TimerSystem.Update() increments this by TimerSystem.Speed each game tick.
	Speed       float64  // Overall update speed of the system; changing this changes how fast the TimerSystem runs. Defaults to 1.
}

// NewTimerSystem creates a new TimerSystem instance.
func NewTimerSystem() *TimerSystem {
	return &TimerSystem{
		Timers: []*Timer{},
		Speed:  1,
	}
}

// AfterTicks creates a new Timer that will elapse after tickCount ticks, running the onElapsed() function when it does so.
// This will happen on whatever thread TimerSystem.Update() is called on (most probably the main thread).
func (ts *TimerSystem) AfterTicks(tickCount TimeUnit, onElapsed func()) *Timer {

	if onElapsed == nil {
		panic("error: onElapsed cannot be nil")
	}

	newTimer := &Timer{
		timerSystem: ts,
		StartTick:   ts.CurrentTime,
		duration:    tickCount,
		OnExecute:   onElapsed,
	}

	ts.Timers = append(ts.Timers, newTimer)

	return newTimer

}

// After creates a new Timer that will elapse after the given duration, running the onElapsed() function when it does so.
// Note that the granularity for conversion from time.Duration is whole ticks, so fractions will be rounded down.
// For example, if your game runs at 60 FPS / TPS, then a tick is 16.67 milliseconds. If you pass a duration of 20 milliseconds,
// the timer will trigger after one tick. If you pass a duration of 16 milliseconds, the timer will trigger immediately.
// This will happen on whatever thread TimerSystem.Update() is called on (most probably the main thread).
func (ts *TimerSystem) After(duration time.Duration, onElapsed func()) *Timer {
	t := ts.AfterTicks(0, onElapsed)
	t.duration = ToTimeUnit(duration)
	return t
}

// Update updates the TimerSystem and triggers any Timers that have elapsed. This should be called once
// per frame in your game's update loop. Note that timers will no longer be accurate if Ebitengine's TPS is changed
// while they are running.
func (ts *TimerSystem) Update() {

	// By looping in reverse, we can freely remove timers while iterating without missing any timers.
	for i := len(ts.Timers) - 1; i >= 0; i-- {

		timer := ts.Timers[i]

		if timer.State == StatePaused {
			timer.StartTick += TimeUnit(ts.Speed)
		} else if timer.State == StateRunning && ts.CurrentTime-timer.StartTick >= timer.duration {

			timer.OnExecute()

			// if it's not looping, we need to remove it from the timers list

			if !timer.Loop {
				timer.State = StateFinished
				ts.removeTimer(timer)
			} else {
				timer.StartTick = ts.CurrentTime
			}

		}

	}

	if ts.Speed < 0 {
		panic("error: speed can't be below 0")
	}

	ts.CurrentTime += TimeUnit(ts.Speed)

}

// remove a timer from the TimerSystem.
func (ts *TimerSystem) removeTimer(timer *Timer) {

	for i, t := range ts.Timers {
		if timer == t {
			ts.Timers[i] = nil
			ts.Timers = append(ts.Timers[:i], ts.Timers[i+1:]...)
		}
	}

}

// Clear cancels all Timers that belong to the TimerSystem and removes them from the TimerSystem. This is
// safe to call from a Timer's elapsing function.
func (ts *TimerSystem) Clear() {

	for _, t := range ts.Timers {
		if t.State != StateFinished {
			t.State = StateCanceled
		}
	}
	ts.Timers = []*Timer{}

}
