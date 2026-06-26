package tuix

// import (
// 	"strings"
// )

// // ─── Window ─────────────────────────────────────────────────────────────────

// type Window struct {
// 	ID    string
// 	Title string
// 	// Content Element
// 	X, Y    int
// 	Width   int
// 	Height  int
// 	Visible bool
// 	OnClose func()
// 	Content func(focused bool, availHeight int) Element
// }

// var (
// 	windows  = make(map[string]*Window)
// 	winStack = []string{}
// )

// // Show opens a window with default size (50x20)
// func Show(id, title string, content func(focused bool, availHeight int) Element) {
// 	ShowAt(id, title, content, 10, 3, 50, 20)
// }

// // ShowAt opens a window at specific position with custom size
// func ShowAt(id, title string, content func(focused bool, availHeight int) Element, x, y, width, height int) {
// 	if _, exists := windows[id]; exists {
// 		Close(id)
// 	}

// 	win := &Window{
// 		ID:      id,
// 		Title:   title,
// 		Content: content,
// 		X:       x,
// 		Y:       y,
// 		Width:   width,
// 		Height:  height,
// 		Visible: true,
// 	}

// 	windows[id] = win
// 	winStack = append(winStack, id)
// }

// // ShowCentered opens a window centered on screen
// func ShowCentered(id, title string, content func(focused bool, availHeight int) Element, width, height int) {
// 	// You'll need to get terminal size here
// 	// For now, using default position
// 	termWidth := 80
// 	termHeight := 24

// 	x := (termWidth - width) / 2
// 	y := (termHeight - height) / 2

// 	ShowAt(id, title, content, x, y, width, height)
// }

// // Close closes a window
// func Close(id string) {
// 	if win, ok := windows[id]; ok {
// 		if win.OnClose != nil {
// 			win.OnClose()
// 		}
// 		win.Visible = false
// 		delete(windows, id)

// 		newStack := []string{}
// 		for _, w := range winStack {
// 			if w != id {
// 				newStack = append(newStack, w)
// 			}
// 		}
// 		winStack = newStack
// 	}
// }

// // CloseCurrent closes the topmost window
// func CloseCurrent() {
// 	if len(winStack) > 0 {
// 		last := winStack[len(winStack)-1]
// 		Close(last)
// 	}
// }

// // IsOpen checks if a window is open
// func IsOpen(id string) bool {
// 	_, ok := windows[id]
// 	return ok
// }

// // HasAnyWindows returns true if any window is open
// func HasAnyWindows() bool {
// 	return len(windows) > 0
// }

// // GetCurrent returns the topmost window
// func GetCurrent() *Window {
// 	if len(winStack) == 0 {
// 		return nil
// 	}
// 	last := winStack[len(winStack)-1]
// 	return windows[last]
// }

// // OnClose sets the close handler
// func OnClose(id string, fn func()) {
// 	if win, ok := windows[id]; ok {
// 		win.OnClose = fn
// 	}
// }

// // ─── Render ─────────────────────────────────────────────────────────────────

// // RenderWindows renders all visible windows as overlays
// // renderWindow builds the chrome for one window
// func renderWindow(win *Window, focused bool) Element {
// 	borderColor := Cyan
// 	if !focused {
// 		borderColor = BrightBlack
// 	}

// 	titleStyle := NewStyle().Bold(true).Foreground(White)
// 	if !focused {
// 		titleStyle = NewStyle().Foreground(BrightBlack)
// 	}

// 	sep := strings.Repeat("─", win.Width-4)

// 	closeHint := Text(" ESC to close", NewStyle().Foreground(BrightBlack))
// 	if !focused {
// 		closeHint = Text("", NewStyle())
// 	}

// 	// Calculate available height for content
// 	// Window height - (title bar 1 + separator 1 + padding 2 + close hint 1)
// 	availableHeight := win.Height - 5

// 	return Box(
// 		Props{
// 			Direction: Column,
// 			Gap:       0,
// 			Padding:   [4]int{1, 2, 1, 2},
// 			Width:     Fixed(win.Width),
// 			Height:    Fixed(win.Height),
// 		},
// 		NewStyle().
// 			Background(Black).
// 			Border(Border{
// 				Top:    true,
// 				Right:  true,
// 				Bottom: true,
// 				Left:   true,
// 				Chars:  BorderRounded,
// 				Color:  borderColor,
// 			}),

// 		// Title bar
// 		Box(
// 			Props{Direction: Row, Padding: [4]int{0, 1, 0, 1}},
// 			NewStyle().Background(Blue),
// 			Text(" "+win.Title+" ", titleStyle),
// 			Text(" [X] ", NewStyle().Foreground(Red)),
// 		),

// 		// Separator
// 		Text(sep, NewStyle().Foreground(BrightBlack)),

// 		// Content receives available height
// 		win.Content(focused, availableHeight),

// 		closeHint,
// 	)
// }

// // func RenderWindows() Element {
// // 	if len(windows) == 0 {
// // 		return Text("", NewStyle())
// // 	}

// // 	children := []Element{}
// // 	for _, id := range winStack {
// // 		win, ok := windows[id]
// // 		if !ok || !win.Visible {
// // 			continue
// // 		}

// // 		// Build window content with fixed width/height
// // 		windowContent := Box(
// // 			Props{
// // 				Direction: Column,
// // 				Gap:       1,
// // 				Padding:   [4]int{1, 2, 1, 2},
// // 				Width:     Fixed(win.Width),
// // 				Height:    Fixed(win.Height),
// // 			},
// // 			NewStyle().
// // 				Background(Black).
// // 				Border(Border{
// // 					Top:    true,
// // 					Right:  true,
// // 					Bottom: true,
// // 					Left:   true,
// // 					Chars:  BorderRounded,
// // 					Color:  Cyan,
// // 				}),
// // 			// Title bar
// // 			Box(
// // 				Props{
// // 					Direction: Row,
// // 					Padding:   [4]int{0, 1, 0, 1},
// // 				},
// // 				NewStyle().Background(Blue),
// // 				Text(" "+win.Title+" ", NewStyle().Bold(true).Foreground(White)),
// // 				Text(" [X] ", NewStyle().Foreground(Red)),
// // 			),
// // 			// Separator
// // 			Text(strings.Repeat("─", win.Width-2), NewStyle().Foreground(BrightBlack)),
// // 			// Content
// // 			win.Content,
// // 			// Close hint
// // 			Text(" ESC to close", NewStyle().Foreground(BrightBlack)),
// // 		)

// // 		// Each window is an Overlay at its position
// // 		children = append(children, Overlay(win.X, win.Y, windowContent))
// // 	}

// // 	if len(children) == 0 {
// // 		return Text("", NewStyle())
// // 	}

// // 	return Box(Props{Direction: Column}, NewStyle(), children...)
// // }

// func RenderWindows() Element {
// 	if len(windows) == 0 {
// 		return Text("", NewStyle())
// 	}

// 	focusedID := winStack[len(winStack)-1]
// 	children := []Element{}

// 	for i, id := range winStack {
// 		win, ok := windows[id]
// 		if !ok || !win.Visible {
// 			continue
// 		}

// 		isFocused := id == focusedID

// 		// Stack windows with offset
// 		x := win.X + i*2
// 		y := win.Y + i

// 		// Pass window and focus state to renderWindow
// 		children = append(children, Overlay(x, y, renderWindow(win, isFocused)))
// 	}

// 	if len(children) == 0 {
// 		return Text("", NewStyle())
// 	}

// 	return Box(Props{Direction: Column}, NewStyle(), children...)
// }

// // RenderAllWindowsAsOverlay renders all windows as a single overlay
// func RenderAllWindowsAsOverlay() Element {
// 	if len(windows) == 0 {
// 		return Text("", NewStyle())
// 	}
// 	return RenderWindows()
// }
