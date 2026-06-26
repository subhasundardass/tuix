// internal/app/screen.go
package app

import "github.com/subhasundardass/tuix/tuix"

// screen.go is the screen stack manager.
// It reads and writes the shared appState declared in state.go.
// All functions here are package-level so app.go can wire them
// into the AppContext via ctx.Set().

// PushScreen adds a screen to the top of the stack and makes it current.
func PushScreen(screenID string) {

	// tuix.Debug("PushScreen called with:", screenID)
	// tuix.Debug("Before push - currentPage:", state.currentPage)
	// tuix.Debug("Before push - screenStack:", state.screenStack)

	state.mu.Lock()
	defer state.mu.Unlock()
	state.screenStack = append(state.screenStack, screenID)
	state.currentPage = screenID

	tuix.ResetComponentState()
}

// PopScreen removes the top screen and returns the new current screen ID.
// If the stack has only one entry it stays — you can never pop the root.
func PopScreen() string {
	state.mu.Lock()
	defer state.mu.Unlock()
	if len(state.screenStack) <= 1 {
		return state.screenStack[0]
	}
	state.screenStack = state.screenStack[:len(state.screenStack)-1]
	top := state.screenStack[len(state.screenStack)-1]
	state.currentPage = top

	tuix.ResetComponentState()
	return top
}

// GetScreenStack returns a copy of the current stack.
// A copy is returned so callers cannot mutate internal state.
func GetScreenStack() []string {
	state.mu.RLock()
	defer state.mu.RUnlock()
	cp := make([]string, len(state.screenStack))
	copy(cp, state.screenStack)
	return cp
}

// GetCurrentScreen returns the ID of the currently active screen.
func GetCurrentScreen() string {
	state.mu.RLock()
	defer state.mu.RUnlock()

	if len(state.screenStack) == 0 {
		return "home"
	}
	return state.screenStack[len(state.screenStack)-1]
}

// ReplaceScreen replaces the top of the stack with a new screen.
// Use this for redirects where you do not want the user to go back.
func ReplaceScreen(screenID string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if len(state.screenStack) == 0 {
		state.screenStack = []string{screenID}
	} else {
		state.screenStack[len(state.screenStack)-1] = screenID
	}

	state.currentPage = screenID
}

// ResetStack clears the stack and sets a single root screen.
// Use this on logout or when navigating to a completely new flow.
func ResetStack(rootScreenID string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.screenStack = []string{rootScreenID}

	state.currentPage = rootScreenID
}

// StackSize returns how many screens are on the stack.
func StackSize() int {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return len(state.screenStack)
}

// CanPop returns true when there is more than one screen on the stack.
// Use this to decide whether to show a back button.
func CanPop() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return len(state.screenStack) > 1
}
