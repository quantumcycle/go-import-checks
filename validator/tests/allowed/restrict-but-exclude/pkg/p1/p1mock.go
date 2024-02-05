package p1

import (
	"github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-but-exclude/internal"
	"github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-but-exclude/pkg/p2"
)

func P1Fn() {
	internal.InternalFn()
	p2.P2Func()
}
