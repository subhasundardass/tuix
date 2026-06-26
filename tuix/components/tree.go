package components

import "github.com/subhasundardass/tuix/tuix"

// ─── Clean Tree Component ─────────────────────────────────────────────────────

type TreeNode struct {
	Label    string
	ID       string
	Children []TreeNode
}

type visibleNode struct {
	id     string
	label  string
	depth  int
	isLeaf bool
	isLast bool
}

func Tree(
	treeID string,
	nodes []TreeNode,
	focused bool,
	onChange func(id string),
) tuix.Element {

	focusedIdx, setFocusedIdx := tuix.UseState(0)
	expanded, setExpanded := tuix.UseStateKeyed(treeID+":expanded", map[string]bool{})

	// Build visible tree
	visible := flattenTree(nodes, expanded)

	if len(visible) == 0 {
		return tuix.Text("(empty)", tuix.NewStyle())
	}

	if focusedIdx >= len(visible) {
		setFocusedIdx(len(visible) - 1)
		focusedIdx = len(visible) - 1
	}

	curr := visible[focusedIdx]

	// KEYBOARD NAVIGATION
	if focused {
		switch tuix.CurrentKey.Code {
		case tuix.KeyUp:
			if focusedIdx > 0 {
				setFocusedIdx(focusedIdx - 1)
			}
		case tuix.KeyDown:
			if focusedIdx < len(visible)-1 {
				setFocusedIdx(focusedIdx + 1)
			}
		case tuix.KeyEnter:
			if curr.isLeaf {
				if onChange != nil {
					onChange(curr.id)
				}
			} else {
				next := make(map[string]bool)
				for k, v := range expanded {
					next[k] = v
				}
				next[curr.id] = !expanded[curr.id]
				setExpanded(next)

				if onChange != nil {
					onChange(curr.id)
				}
			}
		}
	}

	// RENDER TREE LINES
	elems := make([]tuix.Element, 0, len(visible))

	for i, node := range visible {
		isFocused := focused && i == focusedIdx

		// Build tree prefix (indentation with lines)
		prefix := ""
		if node.depth > 0 {
			for j := 0; j < node.depth-1; j++ {
				prefix += "│   "
			}
			if node.isLast {
				prefix += "└── "
			} else {
				prefix += "├── "
			}
		}

		// Build line with separate styles for prefix and content
		var lineElems []tuix.Element

		// Prefix (tree lines) in light gray
		if prefix != "" {
			prefixStyle := tuix.NewStyle().Foreground(tuix.Cyan)
			lineElems = append(lineElems, tuix.Text(prefix, prefixStyle))
		}

		// Icon + Label with focus styling
		contentStyle := tuix.NewStyle()
		if isFocused {
			contentStyle = contentStyle.Background(tuix.Blue).Foreground(tuix.White).Bold(true)
		}
		// lineElems = append(lineElems, tuix.Text(icon+node.label, contentStyle))
		lineElems = append(lineElems, tuix.Text(node.label, contentStyle))

		// Combine as single row element
		elems = append(elems, tuix.Box(
			tuix.Props{Direction: tuix.Row},
			tuix.NewStyle(),
			lineElems...,
		))
	}

	return tuix.Box(
		tuix.Props{Direction: tuix.Column},
		tuix.NewStyle(),
		elems...,
	)
}

// Flatten tree into visible nodes based on expand state
func flattenTree(nodes []TreeNode, expanded map[string]bool) []visibleNode {
	var out []visibleNode
	walkTree(nodes, 0, expanded, &out)
	return out
}

func walkTree(nodes []TreeNode, depth int, expanded map[string]bool, out *[]visibleNode) {
	for i, n := range nodes {
		id := n.ID
		if id == "" {
			id = n.Label
		}

		*out = append(*out, visibleNode{
			id:     id,
			label:  n.Label,
			depth:  depth,
			isLeaf: len(n.Children) == 0,
			isLast: i == len(nodes)-1,
		})

		if expanded[id] && len(n.Children) > 0 {
			walkTree(n.Children, depth+1, expanded, out)
		}
	}
}
