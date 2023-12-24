package run

import (
	"github.com/aFlyBird0/cubox-archiver/core"
	cuboxImpl "github.com/aFlyBird0/cubox-archiver/impl/cubox"
)

func archiverWithDelete(archiver core.Archiver, auth, cookie string) core.Archiver {
	return &archiverAndDelete{
		Archiver: archiver,
		auth:     auth,
		cookie:   cookie,
	}
}

type archiverAndDelete struct {
	core.Archiver
	auth   string
	cookie string
}

func (a *archiverAndDelete) ExistingKeys() (map[string]struct{}, error) {
	return a.Archiver.ExistingKeys()
}

func (a *archiverAndDelete) Operate(item *core.Item) error {
	return core.NewOperatorChain(a.Archiver, cuboxImpl.NewDeleteCuboxOperator(a.auth, a.cookie)).Operate(item)
}
