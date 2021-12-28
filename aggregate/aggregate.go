package aggregate

import (
	"github.com/tehbilly/ges"
)

type Applier interface {
	Apply(event ges.Event) error
}
