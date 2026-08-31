package leftmark

import (
	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
	"github.com/issam-assiyadi/leftmark/adapter/fswalk"
	"github.com/issam-assiyadi/leftmark/application"
)

// New wires a fully-functional Service rooted at root. extraIgnore is a set
// of additional .gitignore-style glob patterns applied on top of root's own
// .gitignore.
func New(root string, extraIgnore ...string) *application.Service {
	return application.NewService(root, fswalk.New(extraIgnore...), fsrw.New())
}
