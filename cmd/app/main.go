package main

import (
	"github.com/subhasundardass/tuix/internal/app"
	"github.com/subhasundardass/tuix/tuix"
)

func main() {

	tuiApp := tuix.NewApp(80, 24)
	tuiApp.Run(app.Root, tuix.Props{})
}
