package p1

import (
	"github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-package/internal"
	"github.com/matdurand/go-import-checks/validator/tests/reject/reject-another-package/pkg/p2"
)

func P1Fn() {
	internal.InternalFn()
	p2.P2Func()
}
