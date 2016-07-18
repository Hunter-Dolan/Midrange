package matcher

import "github.com/Hunter-Dolan/midrange/frame"

// Matcher is a base for matching
type Matcher struct {
	frames  []*frame.Frame
	options *Options
}

// Options configures the matcher
type Options struct {
	*frame.GenerationOptions
	NFFTPower int
}

// NewMatcher creates a new matcher
func NewMatcher(frameGeneratingOptions *Options) *Matcher {
	matcher := Matcher{}
	matcher.options = frameGeneratingOptions

	return &matcher
}
