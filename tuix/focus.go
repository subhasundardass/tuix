// Global focus
package tuix

import "sync"

// FocusManager controls global keyboard focus in the application.
//
// It ensures that only one component is "active" at a time,
// and only the focused component should respond to keyboard input.
//
// It supports:
//   - Single focus ID (active component)
//   - Focus order (Tab / Shift+Tab navigation)
//   - Focus stack (for modal dialogs)
type FocusManager struct {
	mu sync.RWMutex

	// Current focused component ID
	current string

	// Ordered list of focusable component IDs (Tab navigation)
	order []string

	// Stack used for modal / temporary focus overrides
	stack []string
}

// var focusChange = make(chan struct{}, 1)

// NewFocusManager creates a new focus system.
func NewFocusManager() *FocusManager {
	return &FocusManager{
		order: make([]string, 0),
		stack: make([]string, 0),
	}
}

//
// -----------------------------
// BASIC FOCUS CONTROL
// -----------------------------

// SetFocus sets the active focused component.
func (f *FocusManager) SetFocus(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	Debug("SetFocus called with: " + id) // ← add this
	if id == "" {
		f.current = ""
		// trigger re-render

		return
	}

	f.current = id
	Debug("SetFocus current is now: " + f.current)

	// trigger re-render
	// select {
	// case focusChange <- struct{}{}:
	// default:
	// }
}

// Current returns the currently focused component ID.
func (f *FocusManager) Current() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.current
}

// IsFocused checks whether a component is currently focused.
func (f *FocusManager) IsFocused(id string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	Debug("IsFocused check: " + id + " current: " + f.current)
	return f.current == id
}

// global accessor
func (a *App) FocusManager() *FocusManager {
	return a.focus
}

//
// -----------------------------
// TAB NAVIGATION
// -----------------------------

// SetOrder defines the Tab navigation order.
func (f *FocusManager) SetOrder(order []string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.order = order
}

// Next moves focus to next element in order (Tab).
func (f *FocusManager) Next() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.order) == 0 {
		return
	}

	start := f.indexOf(f.current)

	for i := 1; i <= len(f.order); i++ {
		next := (start + i) % len(f.order)
		f.current = f.order[next]
		return
	}
}

// Prev moves focus to previous element (Shift+Tab).
func (f *FocusManager) Prev() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.order) == 0 {
		return
	}

	idx := f.indexOf(f.current)
	prev := (idx - 1 + len(f.order)) % len(f.order)
	f.current = f.order[prev]
}

// indexOf finds index of current focus in order list.
func (f *FocusManager) indexOf(id string) int {
	for i, v := range f.order {
		if v == id {
			return i
		}
	}
	return 0
}

//
// -----------------------------
// MODAL / STACK SUPPORT
// -----------------------------

// PushFocus stores current focus and switches to new one.
// Used when opening modal/dialog.
func (f *FocusManager) PushFocus(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.current != "" {
		f.stack = append(f.stack, f.current)
	}

	f.current = id
}

// PopFocus restores previous focus.
// Used when closing modal/dialog.
func (f *FocusManager) PopFocus() {
	f.mu.Lock()
	defer f.mu.Unlock()

	n := len(f.stack)
	if n == 0 {
		return
	}

	last := f.stack[n-1]
	f.stack = f.stack[:n-1]

	f.current = last
}

//
// -----------------------------
// GLOBAL INSTANCE (SIMPLE MODE)
// -----------------------------

var globalFocus = NewFocusManager()

// CURRENT - returns current focus
func Current() string {
	return globalFocus.Current()
}

// SETORDER - sets Tab navigation order
func SetFocusOrder(order []string) {
	globalFocus.SetOrder(order)
}

// Focus sets global focus.
func Focus(id string) {
	globalFocus.SetFocus(id)
}

// Blur clears focus.
func Blur() {
	globalFocus.SetFocus("")
}

// IsFocused checks global focus.
func IsFocused(id string) bool {
	return globalFocus.IsFocused(id)
}

// FocusNext moves to next focusable item.
func FocusNext() {
	globalFocus.Next()
}

// FocusPrev moves to previous focusable item.
func FocusPrev() {
	globalFocus.Prev()
}

// PushFocus opens modal focus.
func PushFocus(id string) {
	globalFocus.PushFocus(id)
}

// PopFocus closes modal focus.
func PopFocus() {
	globalFocus.PopFocus()
}
