package psichi

import (
	"github.com/m-colson/psi"
)

var _ psi.RouterCreater = Register

func Register(r *psi.Router) error {
	*r = NewMux()
	return nil
}
