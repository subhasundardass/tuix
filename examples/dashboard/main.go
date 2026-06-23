// main.go
package dashboard

import (
	"github.com/subhasundardass/tuix/internal/app"
	"github.com/subhasundardass/tuix/tuix"
)

type DashboardProps struct {
	Title string
}

func DashboardScreen(ctx *app.AppContext, props tuix.Props) tuix.Element {

	appNa := ctx.Config.AppName
	// st:= ctx.NavigateTo("home") // It will call setPage of Bootstrap

	return tuix.Box(
		tuix.Props{},
		tuix.NewStyle(),
		tuix.Text("Dashboard", tuix.NewStyle()),
	)
}

// func main() {
// 	app := tuix.NewApp(80, 24)
// 	app.Run(DashboardScreen, tuix.Props{})
// }
