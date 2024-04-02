//go:build release

package blink

func (v *View) ShowDevTools(devtoolsCallbacks ...func(devtools *View)) {
	// disable for release
}
