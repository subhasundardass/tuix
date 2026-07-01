package window

import (
	"sync"

	"github.com/subhasundardass/tuix/tuix"
)

// WindowManager manages all active windows in the application.
// It maintains the registry, Z-order (stacking), and focus state.
// All operations are thread-safe using a mutex.
type WindowManager struct {
	mu            sync.RWMutex
	windows       map[string]*Window
	stack         []string // Z-order: stack[0]=bottom, last=top
	focused       string   // Currently focused window ID
	renderTrigger func()   // Callback to trigger re-renders in tuix
}

// NewWindowManager creates a new window manager instance.
func NewWindowManager() *WindowManager {
	return &WindowManager{
		windows:       make(map[string]*Window),
		stack:         make([]string, 0),
		focused:       "",
		renderTrigger: nil,
	}
}

// ========================================
// RENDER TRIGGER
// ========================================

// SetRenderTrigger sets the callback function that triggers re-renders.
// This should be called from your main App to connect the window system
// with the tuix render loop.
func (wm *WindowManager) SetRenderTrigger(trigger func()) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.renderTrigger = trigger
}

// triggerRender calls the render trigger if it's set.
// This is called internally when window state changes.
func (wm *WindowManager) triggerRender() {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	if wm.renderTrigger != nil {
		wm.renderTrigger()
	}
}

// ========================================
// WINDOW LIFECYCLE
// ========================================

// AddWindow adds a window to the registry and Z-order stack.
// Called when window.Show() is invoked.
func (wm *WindowManager) AddWindow(w *Window) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.windows[w.ID]; !exists {
		wm.windows[w.ID] = w
		wm.stack = append(wm.stack, w.ID)
		w.focused = true
		wm.focused = w.ID

		// Trigger re-render after adding window
		go wm.triggerRender()
	}
}

// RemoveWindow removes a window from the registry and Z-order.
// Called when window.Close() is invoked.
func (wm *WindowManager) RemoveWindow(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.windows[id]; !exists {
		return
	}

	delete(wm.windows, id)

	// Remove from stack
	for i, wid := range wm.stack {
		if wid == id {
			wm.stack = append(wm.stack[:i], wm.stack[i+1:]...)
			break
		}
	}

	// Clear focus if this was the focused window
	if wm.focused == id {
		wm.focused = ""
		// Focus the topmost remaining window
		if len(wm.stack) > 0 {
			topID := wm.stack[len(wm.stack)-1]
			wm.focused = topID
			if w, exists := wm.windows[topID]; exists {
				w.focused = true
			}
		}
	}

	// Trigger re-render after removing window
	go wm.triggerRender()
}

// GetWindow retrieves a window by ID.
func (wm *WindowManager) GetWindow(id string) *Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.windows[id]
}

// ========================================
// Z-ORDER MANAGEMENT
// ========================================

// BringToFront moves a window to the top of the Z-order (front).
// Also updates focus to this window.
func (wm *WindowManager) BringToFront(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Find window in stack
	idx := -1
	for i, wid := range wm.stack {
		if wid == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return // Window not in stack
	}

	// Remove from current position
	wm.stack = append(wm.stack[:idx], wm.stack[idx+1:]...)

	// Add to top (front)
	wm.stack = append(wm.stack, id)

	// Update focus
	for wid, w := range wm.windows {
		w.focused = (wid == id)
	}
	wm.focused = id

	// Trigger re-render after changing Z-order
	go wm.triggerRender()
}

// GetZOrder returns the current Z-order stack (bottom to top).
// Returns a copy to prevent external modification.
func (wm *WindowManager) GetZOrder() []string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	stack := make([]string, len(wm.stack))
	copy(stack, wm.stack)
	return stack
}

// ========================================
// FOCUS MANAGEMENT
// ========================================

// GetFocused returns the ID of the currently focused window.
func (wm *WindowManager) GetFocused() string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.focused
}

// SetFocus sets which window should receive keyboard input.
func (wm *WindowManager) SetFocus(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.windows[id]; !exists {
		return
	}

	for wid, w := range wm.windows {
		w.focused = (wid == id)
	}
	wm.focused = id
}

// IsFocused checks if a window is currently focused.
func (wm *WindowManager) IsFocused(id string) bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.focused == id
}

// FocusNext cycles focus to the next window in Z-order.
// Useful for keyboard navigation (e.g., Tab key).
func (wm *WindowManager) FocusNext() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if len(wm.stack) == 0 {
		return
	}

	// Find current focus index
	currentIdx := -1
	for i, id := range wm.stack {
		if id == wm.focused {
			currentIdx = i
			break
		}
	}

	// If no focus or last item, wrap to first
	nextIdx := currentIdx + 1
	if nextIdx >= len(wm.stack) {
		nextIdx = 0
	}

	// Set focus to next window
	nextID := wm.stack[nextIdx]
	for wid, w := range wm.windows {
		w.focused = (wid == nextID)
	}
	wm.focused = nextID
}

// ========================================
// QUERY METHODS
// ========================================

// GetAll returns all active windows.
func (wm *WindowManager) GetAll() []*Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	windows := make([]*Window, 0, len(wm.windows))
	for _, w := range wm.windows {
		windows = append(windows, w)
	}
	return windows
}

// GetVisible returns all visible windows in Z-order (bottom to top).
func (wm *WindowManager) GetVisible() []*Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	windows := make([]*Window, 0)
	for _, id := range wm.stack {
		if w, exists := wm.windows[id]; exists && w.visible {
			windows = append(windows, w)
		}
	}
	return windows
}

// Count returns the number of open windows.
func (wm *WindowManager) Count() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return len(wm.windows)
}

// CountVisible returns the number of visible windows.
func (wm *WindowManager) CountVisible() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	count := 0
	for _, w := range wm.windows {
		if w.visible {
			count++
		}
	}
	return count
}

// IsAnyModalOpen returns true if any modal window is currently visible.
func (wm *WindowManager) IsAnyModalOpen() bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	for _, w := range wm.windows {
		if w.visible && w.Modal {
			return true
		}
	}
	return false
}

// GetTopVisibleModal returns the topmost visible modal window (if any).
func (wm *WindowManager) GetTopVisibleModal() *Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	// Walk stack from top to bottom
	for i := len(wm.stack) - 1; i >= 0; i-- {
		if w, exists := wm.windows[wm.stack[i]]; exists && w.visible && w.Modal {
			return w
		}
	}
	return nil
}

// HasWindow checks if a window with the given ID exists.
func (wm *WindowManager) HasWindow(id string) bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	_, exists := wm.windows[id]
	return exists
}

// ========================================
// GLOBAL INSTANCE
// ========================================

var globalManager = NewWindowManager()

// SetRenderTrigger sets the global render trigger function.
// This connects the window system to the tuix render loop.
func SetRenderTrigger(trigger func()) {
	globalManager.SetRenderTrigger(trigger)
}

// Create creates a new window with the given content.
// Window is created but NOT shown - you must call Show() to display it.
// This allows you to configure the window before it's displayed.
//
// Usage:
//
//	win := window.Create(content).
//	    SetTitle("My Window").
//	    SetSize(40, 15).
//	    CenterOnScreen(140, 40).
//	    Show()  // ← Call Show() at the end!
func Create(content interface{}) *Window {
	elem, ok := content.(tuix.Element)
	if !ok {
		panic("window.Create: content must be a tuix.Element")
	}
	return NewWindow(elem)
}

// GetManager returns the global window manager instance.
func GetManager() *WindowManager {
	return globalManager
}

// CloseAll closes all open windows.
func CloseAll() {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	// Clear all windows
	for _, id := range globalManager.stack {
		delete(globalManager.windows, id)
	}
	globalManager.stack = make([]string, 0)
	globalManager.focused = ""

	// Trigger re-render after closing all
	go globalManager.triggerRender()
}

// Count returns the number of open windows globally.
func Count() int {
	return globalManager.Count()
}

// CountVisible returns the number of visible windows globally.
func CountVisible() int {
	return globalManager.CountVisible()
}

// GetFocused returns the globally focused window ID.
func GetFocused() string {
	return globalManager.GetFocused()
}

// IsAnyModalOpen returns true if any modal is currently open.
func IsAnyModalOpen() bool {
	return globalManager.IsAnyModalOpen()
}

// GetVisible returns all visible windows globally.
func GetVisible() []*Window {
	return globalManager.GetVisible()
}
