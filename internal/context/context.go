package context

import "github.com/subhasundardass/tuix/tuix"

var DefaultContext = tuix.CreateContext[*AppContext](nil)

// Use returns the current context
func Use() *AppContext {
	return tuix.UseContext(DefaultContext)
}

type AppContext struct {
	currentPage   string
	appName       string
	darkMode      bool
	userName      string
	navigateTo    func(string)
	toggleDark    func()
	pushScreen    func(string)
	replaceScreen func(string)
	popScreen     func() string
	getStack      func() []string
	getCurrent    func() string
}

type AppContextValues struct {
	CurrentPage   string
	AppName       string
	DarkMode      bool
	UserName      string
	NavigateTo    func(string)
	ToggleDark    func()
	PushScreen    func(string)
	ReplaceScreen func(string)
	PopScreen     func() string
	GetStack      func() []string
	GetCurrent    func() string
}

func (c *AppContext) Set(v AppContextValues) {
	// ⭐ Add nil check
	if c == nil {
		return
	}
	c.currentPage = v.CurrentPage
	c.appName = v.AppName
	c.darkMode = v.DarkMode
	c.userName = v.UserName
	c.navigateTo = v.NavigateTo
	c.toggleDark = v.ToggleDark
	c.pushScreen = v.PushScreen
	c.replaceScreen = v.ReplaceScreen
	c.popScreen = v.PopScreen
	c.getStack = v.GetStack
	c.getCurrent = v.GetCurrent
}

func (c *AppContext) Page() string {
	if c == nil {
		return "home"
	}
	return c.currentPage
}

func (c *AppContext) AppName() string {
	if c == nil {
		return "App"
	}
	return c.appName
}

func (c *AppContext) IsDarkMode() bool {
	if c == nil {
		return false
	}
	return c.darkMode
}

func (c *AppContext) UserName() string {
	if c == nil {
		return "Guest"
	}
	return c.userName
}

func (c *AppContext) NavigateTo(page string) {
	if c == nil {
		return
	}
	c.navigateTo(page)
}

func (c *AppContext) ToggleDark() {
	if c == nil || c.toggleDark == nil {
		return
	}
	c.toggleDark()
}

func (c *AppContext) PushScreen(screenID string) {

	if c == nil {
		return
	}
	c.pushScreen(screenID)
}

func (c *AppContext) ReplaceScreen(screenID string) {
	if c == nil || c.replaceScreen == nil {
		return
	}
	c.replaceScreen(screenID)
}

func (c *AppContext) PopScreen() string {
	if c == nil || c.popScreen == nil {
		return "home"
	}
	return c.popScreen()
}

func (c *AppContext) GetStack() []string {
	if c == nil || c.getStack == nil {
		return []string{"home"}
	}
	return c.getStack()
}

func (c *AppContext) GetCurrent() string {
	if c == nil || c.getCurrent == nil {
		return "home"
	}
	return c.getCurrent()
}
