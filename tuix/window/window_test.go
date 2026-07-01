package window

import (
	"testing"

	"github.com/subhasundardass/tuix/tuix"
)

// TestNewWindow tests the creation of a new window
func TestNewWindow(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	if win == nil {
		t.Error("NewWindow returned nil")
	}

	if win.Title != "Window" {
		t.Errorf("Expected default title 'Window', got '%s'", win.Title)
	}

	if win.Width != 40 {
		t.Errorf("Expected default width 40, got %d", win.Width)
	}

	if win.Height != 15 {
		t.Errorf("Expected default height 15, got %d", win.Height)
	}

	if win.visible {
		t.Error("New window should be hidden by default")
	}

	if win.focused {
		t.Error("New window should not be focused by default")
	}
}

// TestWindowSetTitle tests the SetTitle method
func TestWindowSetTitle(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	win.SetTitle("My Test Window")
	if win.Title != "My Test Window" {
		t.Errorf("Expected title 'My Test Window', got '%s'", win.Title)
	}

	// Test fluent interface
	win.SetTitle("Another Title").SetSize(50, 20)
	if win.Title != "Another Title" {
		t.Errorf("Expected title 'Another Title', got '%s'", win.Title)
	}
	if win.Width != 50 || win.Height != 20 {
		t.Errorf("Expected size 50x20, got %dx%d", win.Width, win.Height)
	}
}

// TestWindowSetSize tests the SetSize method
func TestWindowSetSize(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	win.SetSize(60, 25)
	if win.Width != 60 {
		t.Errorf("Expected width 60, got %d", win.Width)
	}
	if win.Height != 25 {
		t.Errorf("Expected height 25, got %d", win.Height)
	}
}

// TestWindowSetPosition tests the SetPosition method
func TestWindowSetPosition(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	win.SetPosition(10, 20)
	if win.X != 10 {
		t.Errorf("Expected X 10, got %d", win.X)
	}
	if win.Y != 20 {
		t.Errorf("Expected Y 20, got %d", win.Y)
	}
}

// TestWindowSetModal tests the SetModal method
func TestWindowSetModal(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	if win.Modal {
		t.Error("New window should not be modal by default")
	}

	win.SetModal(true)
	if !win.Modal {
		t.Error("Window should be modal after SetModal(true)")
	}

	win.SetModal(false)
	if win.Modal {
		t.Error("Window should not be modal after SetModal(false)")
	}
}

// TestWindowShow tests the Show method
func TestWindowShow(t *testing.T) {
	// Clear any existing windows
	CloseAll()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	if win.IsVisible() {
		t.Error("Window should be hidden before Show()")
	}

	win.Show()

	if !win.IsVisible() {
		t.Error("Window should be visible after Show()")
	}

	// Check that window was added to manager
	if Count() != 1 {
		t.Errorf("Expected 1 window in manager, got %d", Count())
	}

	// Cleanup
	CloseAll()
}

// TestWindowHide tests the Hide method
func TestWindowHide(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.Show()

	if !win.IsVisible() {
		t.Error("Window should be visible after Show()")
	}

	win.Hide()

	if win.IsVisible() {
		t.Error("Window should be hidden after Hide()")
	}

	// Window should still be in manager
	if Count() != 1 {
		t.Errorf("Expected 1 window in manager after hide, got %d", Count())
	}

	ClearManager()
}

// TestWindowClose tests the Close method
func TestWindowClose(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.Show()

	if !win.IsVisible() {
		t.Error("Window should be visible after Show()")
	}

	win.Close()

	if win.IsVisible() {
		t.Error("Window should be hidden after Close()")
	}

	// Window should be removed from manager
	if Count() != 0 {
		t.Errorf("Expected 0 windows in manager after close, got %d", Count())
	}

	ClearManager()
}

// TestWindowFocus tests the Focus method
func TestWindowFocus(t *testing.T) {
	ClearManager()

	content1 := tuix.Text("Window 1", tuix.NewStyle())
	content2 := tuix.Text("Window 2", tuix.NewStyle())

	win1 := NewWindow(content1)
	win2 := NewWindow(content2)

	win1.Show()
	win2.Show()

	// win2 should be focused (last shown)
	if !win2.IsFocused() {
		t.Error("Last shown window should be focused")
	}

	// Focus win1
	win1.Focus()

	if !win1.IsFocused() {
		t.Error("win1 should be focused after Focus()")
	}

	if win2.IsFocused() {
		t.Error("win2 should not be focused after win1.Focus()")
	}

	ClearManager()
}

// TestWindowCenterOnScreen tests the CenterOnScreen method
func TestWindowCenterOnScreen(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	win.SetSize(40, 15)
	win.CenterOnScreen(140, 40)

	expectedX := (140 - 40) / 2
	expectedY := (40 - 15) / 2

	if win.X != expectedX {
		t.Errorf("Expected X %d, got %d", expectedX, win.X)
	}
	if win.Y != expectedY {
		t.Errorf("Expected Y %d, got %d", expectedY, win.Y)
	}

	// Test with window larger than screen - it should clamp to 0
	win.SetSize(200, 50)
	win.CenterOnScreen(140, 40)

	// When window is larger than screen, X and Y should be clamped to 0
	if win.X != 0 {
		t.Errorf("Expected X to be clamped to 0 when window larger than screen, got %d", win.X)
	}
	if win.Y != 0 {
		t.Errorf("Expected Y to be clamped to 0 when window larger than screen, got %d", win.Y)
	}
}

// TestWindowGetBounds tests the GetBounds method
func TestWindowGetBounds(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	win.SetPosition(10, 20)
	win.SetSize(50, 25)

	bounds := win.GetBounds()
	expected := [4]int{10, 20, 50, 25}

	if bounds != expected {
		t.Errorf("Expected bounds %v, got %v", expected, bounds)
	}
}

// TestWindowRender tests the Render method
func TestWindowRender(t *testing.T) {
	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	rendered := win.Render()
	if rendered.Text != "Test Content" {
		t.Errorf("Expected rendered text 'Test Content', got '%s'", rendered.Text)
	}
}

// TestWindowManagerAddWindow tests adding windows to manager
func TestWindowManagerAddWindow(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)

	mgr := GetManager()
	mgr.AddWindow(win)

	if Count() != 1 {
		t.Errorf("Expected 1 window, got %d", Count())
	}

	// Test duplicate add
	mgr.AddWindow(win)
	if Count() != 1 {
		t.Errorf("Expected still 1 window after duplicate add, got %d", Count())
	}

	ClearManager()
}

// TestWindowManagerRemoveWindow tests removing windows from manager
func TestWindowManagerRemoveWindow(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.Show()

	if Count() != 1 {
		t.Errorf("Expected 1 window, got %d", Count())
	}

	mgr := GetManager()
	mgr.RemoveWindow(win.ID)

	if Count() != 0 {
		t.Errorf("Expected 0 windows, got %d", Count())
	}

	// Test removing non-existent window
	mgr.RemoveWindow("non-existent-id") // Should not panic
}

// TestWindowManagerBringToFront tests Z-order management
func TestWindowManagerBringToFront(t *testing.T) {
	ClearManager()

	content1 := tuix.Text("Window 1", tuix.NewStyle())
	content2 := tuix.Text("Window 2", tuix.NewStyle())
	content3 := tuix.Text("Window 3", tuix.NewStyle())

	win1 := NewWindow(content1)
	win2 := NewWindow(content2)
	win3 := NewWindow(content3)

	win1.Show()
	win2.Show()
	win3.Show()

	mgr := GetManager()
	zorder := mgr.GetZOrder()

	// Last added should be on top
	if zorder[len(zorder)-1] != win3.ID {
		t.Errorf("Expected win3 on top, got %s", zorder[len(zorder)-1])
	}

	// Bring win1 to front
	mgr.BringToFront(win1.ID)
	zorder = mgr.GetZOrder()

	if zorder[len(zorder)-1] != win1.ID {
		t.Errorf("Expected win1 on top after BringToFront, got %s", zorder[len(zorder)-1])
	}

	ClearManager()
}

// TestWindowManagerGetFocused tests focus management
func TestWindowManagerGetFocused(t *testing.T) {
	ClearManager()

	content1 := tuix.Text("Window 1", tuix.NewStyle())
	content2 := tuix.Text("Window 2", tuix.NewStyle())

	win1 := NewWindow(content1)
	win2 := NewWindow(content2)

	win1.Show()
	win2.Show()

	mgr := GetManager()
	focused := mgr.GetFocused()

	if focused != win2.ID {
		t.Errorf("Expected focused window %s, got %s", win2.ID, focused)
	}

	// Focus win1
	mgr.SetFocus(win1.ID)
	focused = mgr.GetFocused()

	if focused != win1.ID {
		t.Errorf("Expected focused window %s after SetFocus, got %s", win1.ID, focused)
	}

	ClearManager()
}

// TestWindowManagerIsAnyModalOpen tests modal detection
func TestWindowManagerIsAnyModalOpen(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.SetModal(true)
	win.Show()

	mgr := GetManager()

	if !mgr.IsAnyModalOpen() {
		t.Error("IsAnyModalOpen should return true when modal window is open")
	}

	win.Close()

	if mgr.IsAnyModalOpen() {
		t.Error("IsAnyModalOpen should return false after modal window is closed")
	}

	ClearManager()
}

// TestWindowManagerGetTopVisibleModal tests getting top modal
func TestWindowManagerGetTopVisibleModal(t *testing.T) {
	ClearManager()

	content1 := tuix.Text("Modal 1", tuix.NewStyle())
	content2 := tuix.Text("Modal 2", tuix.NewStyle())

	win1 := NewWindow(content1)
	win1.SetModal(true)
	win1.Show()

	win2 := NewWindow(content2)
	win2.SetModal(true)
	win2.Show()

	mgr := GetManager()
	topModal := mgr.GetTopVisibleModal()

	if topModal.ID != win2.ID {
		t.Errorf("Expected top modal %s, got %s", win2.ID, topModal.ID)
	}

	ClearManager()
}

// TestWindowManagerCount tests Count method
func TestWindowManagerCount(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())

	if Count() != 0 {
		t.Errorf("Expected 0 windows, got %d", Count())
	}

	win1 := NewWindow(content)
	win1.Show()

	if Count() != 1 {
		t.Errorf("Expected 1 window, got %d", Count())
	}

	win2 := NewWindow(content)
	win2.Show()

	if Count() != 2 {
		t.Errorf("Expected 2 windows, got %d", Count())
	}

	win1.Close()

	if Count() != 1 {
		t.Errorf("Expected 1 window after close, got %d", Count())
	}

	ClearManager()
}

// TestCloseAll tests closing all windows
func TestCloseAll(t *testing.T) {
	ClearManager()

	content := tuix.Text("Test Content", tuix.NewStyle())

	win1 := NewWindow(content)
	win2 := NewWindow(content)
	win3 := NewWindow(content)

	win1.Show()
	win2.Show()
	win3.Show()

	if Count() != 3 {
		t.Errorf("Expected 3 windows, got %d", Count())
	}

	CloseAll()

	if Count() != 0 {
		t.Errorf("Expected 0 windows after CloseAll, got %d", Count())
	}
}

// TestWindowString tests the String method
func TestWindowString(t *testing.T) {
	// Reset the counter for predictable testing
	windowCounter = 0

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.SetTitle("Test Window")
	win.SetSize(50, 20)
	win.SetPosition(10, 5)

	str := win.String()
	expected := "Window{ID: win-1, Title: Test Window, Size: 50x20, Pos: (10,5), Visible: false, Modal: false}"

	if str != expected {
		t.Errorf("Expected string %s, got %s", expected, str)
	}
}

// TestWindowClone tests the Clone method
func TestWindowClone(t *testing.T) {
	// Reset counter for predictable testing
	windowCounter = 0

	content := tuix.Text("Test Content", tuix.NewStyle())
	win := NewWindow(content)
	win.SetTitle("Original Window")
	win.SetSize(50, 20)
	win.SetPosition(10, 5)
	win.SetModal(true)

	clone := win.Clone()

	if clone.Title != "Original Window (Copy)" {
		t.Errorf("Expected clone title 'Original Window (Copy)', got '%s'", clone.Title)
	}

	if clone.Width != 50 || clone.Height != 20 {
		t.Errorf("Expected clone size 50x20, got %dx%d", clone.Width, clone.Height)
	}

	if clone.X != 15 || clone.Y != 10 {
		t.Errorf("Expected clone position (15,10), got (%d,%d)", clone.X, clone.Y)
	}

	if clone.ID == win.ID {
		t.Error("Clone should have different ID")
	}

	if clone.visible {
		t.Error("Clone should be hidden by default")
	}
}

// Helper function to clear manager between tests
func ClearManager() {
	CloseAll()
	// Reset the global manager
	globalManager = NewWindowManager()
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	// Setup
	ClearManager()
	windowCounter = 0

	// Run tests
	m.Run()

	// Cleanup
	ClearManager()
}
