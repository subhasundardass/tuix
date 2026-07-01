package tuix

import "reflect"

type Effect struct {
	fn      func() func()
	deps    []any
	cleanup func()
	dirty   bool
}

var State []any
var StateCursor int = 0

var Effects []Effect
var EffectCursor int = 0

// KeyedState holds state for UseStateKeyed — a map keyed by a stable
// string identity rather than cursor position. This lets components like
// Tree maintain independent expand/collapse state per node across renders
// where the number of visible nodes changes (collapsed subtrees disappear,
// which would shift the positional StateCursor for every node after them).
var KeyedState = map[string]any{}

var pendingRender bool

func UseState[T any](initial T) (T, func(T)) {
	idx := StateCursor
	StateCursor++

	if idx >= len(State) {
		State = append(State, initial)
	}

	current := State[idx].(T)

	setter := func(next T) {
		State[idx] = next
		pendingRender = true
	}

	return current, setter
}

// UseStateKeyed is like UseState but keyed by a stable string instead of
// cursor position. Use this whenever the number of hook calls in a render
// can vary — e.g. a tree node that may or may not render children depending
// on expand state. Positional UseState would corrupt all subsequent state
// slots when nodes are toggled; UseStateKeyed is immune because it always
// looks up by key regardless of render order.
//
// key must be globally unique across your entire component tree for the
// lifetime of the app. For tree nodes, combine the node's ID with its
// depth or path: "tree-node-/src/main.go" or "tree-node-0-1-2".
func UseStateKeyed[T any](key string, initial T) (T, func(T)) {
	if _, exists := KeyedState[key]; !exists {
		KeyedState[key] = initial
	}

	current := KeyedState[key].(T)

	setter := func(next T) {
		KeyedState[key] = next
		pendingRender = true
	}

	return current, setter
}

// UseEffect registers a side-effect to run after paint whenever its
// deps change. fn returns an optional cleanup called before the next
// run or when the component unmounts.
//
// Deps comparison uses reflect.DeepEqual so slice/map/struct deps are
// compared by value, not pointer — a plain != on interface{} would
// panic at runtime for any non-comparable dep type (slice, map, etc).
func UseEffect(fn func() func(), deps []any) {
	idx := EffectCursor
	EffectCursor++

	newEffect := Effect{fn: fn, deps: deps, dirty: true}

	if idx >= len(Effects) {
		// ⭐ New effect - append
		Effects = append(Effects, newEffect)
		return
	}

	// ⭐ Existing effect - check if deps changed
	existing := &Effects[idx]
	changed := len(existing.deps) != len(newEffect.deps)
	if !changed {
		for i, dep := range newEffect.deps {
			if !reflect.DeepEqual(existing.deps[i], dep) {
				changed = true
				break
			}
		}
	}

	if changed {
		// ⭐ Keep the old cleanup - RunEffects will call it
		existing.dirty = true
		existing.fn = newEffect.fn
		existing.deps = newEffect.deps
	} else {
		existing.dirty = false
	}
}

// RunEffects runs all effects marked dirty since the last render.
// Called after screen flush so effects fire after paint, matching
// React semantics.
func RunEffects() {
	for i := range Effects {
		if !Effects[i].dirty {
			continue
		}
		// ⭐ Call old cleanup if it exists
		if Effects[i].cleanup != nil {
			Effects[i].cleanup()
		}
		// ⭐ Run the new effect and store its cleanup
		Effects[i].cleanup = Effects[i].fn()
		Effects[i].dirty = false
	}
}

// Context carries a value down the component tree without prop-drilling.
// The zero value is not usable — construct with CreateContext.
type Context[T any] struct {
	defaultValue T
	stack        []T
}

// CreateContext returns a new Context whose UseContext readers see
// defaultValue when no enclosing Provide is active.
func CreateContext[T any](defaultValue T) *Context[T] {
	return &Context[T]{defaultValue: defaultValue}
}

// Provide pushes value onto the context's stack, runs render (during
// which any descendant calling UseContext observes value), then pops
// via defer so a panic in render still unwinds the stack cleanly.
func (c *Context[T]) Provide(value T, render func() Element) Element {
	c.stack = append(c.stack, value)
	defer func() { c.stack = c.stack[:len(c.stack)-1] }()
	return render()
}

// UseContext returns the innermost active Provide value, or
// defaultValue if no Provide is currently on the stack.
func UseContext[T any](c *Context[T]) T {
	if len(c.stack) == 0 {
		return c.defaultValue
	}
	return c.stack[len(c.stack)-1]
}

// ResetComponentState clears all positional state slots.
// Call this whenever the active screen changes so stale state
// from the previous screen does not corrupt the new screen's slots.
func ResetComponentState() {
	State = nil
	StateCursor = 0
	Effects = nil
	EffectCursor = 0
}
