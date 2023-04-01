package filter

import (
	"github.com/aFlyBird0/cubox-archiver/cubox"
)

type Filter interface {
	Remain(item *cubox.Item) bool
}
