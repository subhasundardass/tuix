package context

import "github.com/subhasundardass/tuix/tuix"

var DefaultContext = tuix.CreateContext[*AppContext](nil)

// Use — call once at top of any component, like fiber.Ctx
func Use() *AppContext {
	return tuix.UseContext(DefaultContext)
}

type AppContext struct {
	currentPage string
	appName     string
	darkMode    bool
	userName    string
	navigateTo  func(page string)
	toggleDark  func()

	//Screen stack for navigation history
	screenStack []string
	pushScreen  func(string)
	popScreen   func() string
	getStack    func() []string
	getCurrent  func() string
}

type AppContextValues struct {
	CurrentPage string
	AppName     string
	DarkMode    bool
	UserName    string
	NavigateTo  func(string)
	ToggleDark  func()

	// Screen stack functions
	PushScreen    func(string)
	ReplaceScreen func(string)
	PopScreen     func() string
	GetStack      func() []string
	GetCurrent    func() string
}

func (c *AppContext) Set(v AppContextValues) {
	c.currentPage = v.CurrentPage
	c.appName = v.AppName
	c.darkMode = v.DarkMode
	c.userName = v.UserName
	c.navigateTo = v.NavigateTo
	c.toggleDark = v.ToggleDark

	// Screen stack
	c.pushScreen = v.PushScreen
	c.popScreen = v.PopScreen
	c.getStack = v.GetStack
	c.getCurrent = v.GetCurrent
}

func (c *AppContext) Page() string     { return c.currentPage }
func (c *AppContext) AppName() string  { return c.appName }
func (c *AppContext) IsDarkMode() bool { return c.darkMode }
func (c *AppContext) UserName() string { return c.userName }

func (c *AppContext) ToggleDark() { c.toggleDark() }

// ─── Screen Stack Methods ───────────────────────────────────────────────

// PushScreen adds a screen to the stack and navigates to it
func (c *AppContext) PushScreen(screenID string) {
	if c.pushScreen != nil {
		c.pushScreen(screenID)
	}
}

// PopScreen removes the top screen from the stack
// Returns the new current screen ID
func (c *AppContext) PopScreen() string {
	if c.popScreen != nil {
		return c.popScreen()
	}
	return "home"
}

// GetStack returns a copy of the screen stack
func (c *AppContext) GetStack() []string {
	if c.getStack != nil {
		return c.getStack()
	}
	return []string{"home"}
}

// GetCurrent returns the current screen ID
func (c *AppContext) GetCurrent() string {
	if c.getCurrent != nil {
		return c.getCurrent()
	}
	return c.currentPage
}

// StackSize returns the number of screens in the stack
func (c *AppContext) StackSize() int {
	return len(c.GetStack())
}
