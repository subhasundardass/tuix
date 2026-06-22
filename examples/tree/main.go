package main

import (
	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

var sidebarTree = []components.TreeNode{
	{
		ID:    "cmd",
		Label: "cmd",
		Children: []components.TreeNode{
			{
				ID:    "cmd-app",
				Label: "app",
				Children: []components.TreeNode{
					{ID: "cmd-app-main", Label: "main.go"},
					{ID: "cmd-app-bootstrap", Label: "bootstrap.go"},
				},
			},
		},
	},
	{
		ID:    "internal",
		Label: "internal",
		Children: []components.TreeNode{
			{
				ID:    "internal-ui",
				Label: "ui",
				Children: []components.TreeNode{
					{
						ID:    "internal-ui-components",
						Label: "components",
						Children: []components.TreeNode{
							{ID: "component-header", Label: "header.go"},
							{ID: "component-footer", Label: "footer.go"},
							{ID: "component-sidebar", Label: "sidebar.go"},
							{ID: "component-tree", Label: "tree.go"},
							{ID: "component-table", Label: "table.go"},
							{ID: "component-modal", Label: "modal.go"},
						},
					},
					{
						ID:    "internal-ui-screens",
						Label: "screens",
						Children: []components.TreeNode{
							{ID: "screen-dashboard", Label: "dashboard.go"},
							{ID: "screen-users", Label: "users.go"},
							{ID: "screen-settings", Label: "settings.go"},
						},
					},
					{
						ID:    "internal-ui-layouts",
						Label: "layouts",
						Children: []components.TreeNode{
							{ID: "layout-master", Label: "master.go"},
						},
					},
				},
			},
			{
				ID:    "internal-controller",
				Label: "controller",
				Children: []components.TreeNode{
					{ID: "controller-dashboard", Label: "dashboard.go"},
					{ID: "controller-users", Label: "users.go"},
				},
			},
			{
				ID:    "internal-navigation",
				Label: "navigation",
				Children: []components.TreeNode{
					{ID: "nav-router", Label: "router.go"},
					{ID: "nav-screenmanager", Label: "screen_manager.go"},
				},
			},
			{
				ID:    "internal-state",
				Label: "state",
				Children: []components.TreeNode{
					{ID: "state-app", Label: "app_state.go"},
				},
			},
		},
	},
	{
		ID:    "pkg",
		Label: "pkg",
		Children: []components.TreeNode{
			{ID: "pkg-config", Label: "config"},
			{ID: "pkg-utils", Label: "utils"},
		},
	},
	{
		ID:    "docs",
		Label: "docs",
		Children: []components.TreeNode{
			{ID: "docs-readme", Label: "README.md"},
			{ID: "docs-api", Label: "API.md"},
		},
	},
	{
		ID:    "gomod",
		Label: "go.mod",
	},
	{
		ID:    "gosum",
		Label: "go.sum",
	},
}

func App(props tuix.Props) tuix.Element {
	selected, setSelected := tuix.UseState("")

	titleStyle := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dimStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)
	selectedStyle := tuix.NewStyle().Foreground(tuix.BrightYellow)
	borderStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)

	selectedLabel := tuix.If(
		selected != "",
		tuix.Text("selected: "+selected, selectedStyle),
		tuix.Text("navigate with ↑/↓ · Enter to expand/select", dimStyle),
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Row, Padding: [4]int{1, 1, 1, 1}},
		tuix.NewStyle(),

		// Sidebar tree panel
		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),
			tuix.Text("◆ Files", titleStyle),
			tuix.Text("─────────────", borderStyle),

			components.Tree(
				"sidebar",
				sidebarTree,
				true,
				func(id string) {
					setSelected(id)
				},
			),
		),

		// Vertical divider
		tuix.Text(" │ ", borderStyle),

		// Main panel
		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("◆ Detail", titleStyle),
			selectedLabel,
			tuix.Text("ctrl-c to quit", dimStyle),
		),
	)
}

func main() {
	app := tuix.NewApp(70, 20)
	app.Run(App, tuix.Props{})
}
