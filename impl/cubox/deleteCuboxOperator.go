package cubox

import (
	"github.com/parnurzeal/gorequest"
	"go.uber.org/multierr"

	"github.com/aFlyBird0/cubox-archiver/core"
)

type DeleteCuboxOperator struct {
	auth   string
	cookie string
}

func NewDeleteCuboxOperator(auth, cookie string) *DeleteCuboxOperator {
	return &DeleteCuboxOperator{auth: auth, cookie: cookie}
}

func (o *DeleteCuboxOperator) Operate(item *core.Item) error {
	url := "https://cubox.pro/c/api/search_engine/delete/" + item.UserSearchEngineID

	req := gorequest.New().Post(url)
	req = SetGoRequestHeader(req, o.auth, o.cookie)
	_, _, errs := req.End()

	return multierr.Combine(errs...)

}
