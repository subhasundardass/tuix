package layout

import (
	"github.com/subhasundardass/tuix/tuix"
	"github.com/subhasundardass/tuix/tuix/components"
)

var sidebarTree = []components.TreeNode{
	{
		ID:    "dashboard",
		Label: "Dashboard",
		Children: []components.TreeNode{
			{
				ID:    "cmd-app",
				Label: "app",
				Children: []components.TreeNode{
					{ID: "home", Label: "Home"},
					{ID: "settings", Label: "Settings"},
					{ID: "about", Label: "About"},
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

func SidebarTree(props tuix.Props) tuix.Element {
	// selected, setSelected := tuix.UseState("")

	// ctx := context.Use()

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Padding: [4]int{1, 0, 0, 1}, Width: tuix.Fixed(30), Gap: 0},
		tuix.NewStyle().Border(tuix.Border{
			Top: true, Right: true, Bottom: true, Left: true,
			Chars: tuix.BorderRounded, Color: tuix.BrightBlack,
			Title: "Navigation",
		}),

		// Sidebar tree panel
		tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 0},
			tuix.NewStyle(),

			components.Tree(
				"sidebar",
				sidebarTree,
				true,
				func(id string) {
					tuix.Debug("Navigating to --> ", id)
					// ctx.PushScreen(id)
					// tuix.Show("about", "About Us", 60, 20, func(focused bool) tuix.Element {
					// 	return screen.AboutPage(ctx, tuix.Props{}, focused)
					// })
				},
			),
		),
	)
}
