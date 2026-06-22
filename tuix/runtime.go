package tuix

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	screen   *Screen
	renderer *ComponentRenderer
}

// NewApp creates and initializes a new terminal application.
//
// NewApp performs the complete application bootstrap process:
//
//  1. Creates a Screen that writes to os.Stdout.
//  2. Puts the terminal into raw mode.
//  3. Hides the cursor and enables bracketed paste mode.
//  4. Detects the actual terminal dimensions.
//  5. Creates the root Renderer.
//  6. Returns a fully initialized App instance.
//
// The width and height arguments act as fallback dimensions when the
// terminal size cannot be determined (for example when stdout is
// redirected or running in certain CI environments).
//
// In normal interactive terminals, the application's dimensions are
// automatically expanded to the real terminal viewport, regardless of
// the values passed to this constructor.
//
// Example:
//
//	app := tuix.NewApp(80, 24)
//	defer app.Stop()
//
//	app.Render(func() tuix.Element {
//	    return tuix.Text("Hello, TUIX!")
//	})
//
//	app.Run()
//
// Example using automatic terminal sizing:
//
//	app := tuix.NewApp(0, 0)
//	defer app.Stop()
//
//	app.Run()
func NewApp(width, height int) *App {

	screen := NewScreenWriter(width, height, os.Stdout)
	screen.Start()

	// Prefer the real terminal dimensions over the constructor args so
	// layout fills the actual viewport. The args remain a fallback for
	// environments where term.GetSize fails (e.g. piped output).
	if screen.termCols > 0 && screen.termRows > 0 {
		screen.SetDimensions(screen.termCols, screen.termRows)
	} else {
		screen.SetDimensions(width, height)
	}

	renderer := NewRenderer(screen)

	return &App{
		screen:   screen,
		renderer: renderer,
	}
}

var ticker = make(chan bool, 1)
var CurrentTick bool = false

var exitCh = make(chan struct{}, 1)

// Exit requests the running application to stop gracefully.
func Exit() {
	select {
	case exitCh <- struct{}{}:
	default:
	}
}

func (a *App) Run(fn func(props Props) Element, props Props) {
	a.Render(fn, props)

	select {
	case <-exitCh:
	default:
	}

	quit := make(chan struct{})
	var quitOnce sync.Once
	requestQuit := func() {
		quitOnce.Do(func() { close(quit) })
	}

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)

	go func() {
		tick := false
		for {
			time.Sleep(time.Millisecond * 500)
			tick = !tick
			ticker <- tick
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		var scanner KeyScanner
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				requestQuit()
				return
			}
			for _, key := range scanner.Feed(buf[:n]) {
				if key.Code == KeyCtrlC {
					requestQuit()
					return
				}
				Keys <- key
			}
		}
	}()

	for {
		select {
		case <-quit:
			a.screen.Stop()
			return
		case <-exitCh:
			requestQuit()
		case key := <-Keys:
			CurrentKey = key
			a.Render(fn, props)
		case tick := <-ticker:
			CurrentTick = tick
			a.Render(fn, props)
		case <-resize:
			a.screen.HandleResize()
			a.Render(fn, props)
		}
	}
}

func (a *App) Render(fn func(props Props) Element, props Props) {
	// Pass 1: process key events and mutate state
	StateCursor = 0
	EffectCursor = 0
	fn(props)

	// Pass 2: render with updated state; key is now consumed
	CurrentKey = Key{}
	StateCursor = 0
	EffectCursor = 0
	next := fn(props)

	a.renderer.Render(next)
	a.screen.Flush()

	pendingRender = false
	RunEffects()

	if pendingRender {
		CurrentKey = Key{}
		StateCursor = 0
		EffectCursor = 0
		next := fn(props)

		a.renderer.Render(next)
		a.screen.Flush()
	}
}
