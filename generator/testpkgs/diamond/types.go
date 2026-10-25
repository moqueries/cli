// Package diamond contains interfaces that embed a common interface through
// multiple paths
package diamond

// Base is embedded by both Left and Right
type Base interface {
	Base(sParam string) error
}

// Left embeds Base
type Left interface {
	Base
	Left()
}

// Right embeds Base
type Right interface {
	Base
	Right()
}

// Diamond embeds Base twice, once via Left and once via Right
type Diamond interface {
	Left
	Right
}
