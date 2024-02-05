package p1

import (
	"github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-same-package/internal"
	"github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-same-package/pkg/p2"
)

func P1Fn() {
	internal.InternalFn()
	p2.P2Func()
}
