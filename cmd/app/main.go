package main

import (
	"github.com/subhasundardass/tuix/internal/app"
	"github.com/subhasundardass/tuix/tuix"
)

func main() {

	tuiApp := tuix.NewApp(0, 0) // Full Screen

	tuiApp.Run(app.Root, tuix.Props{})
}
