package api

import (
	d1 "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-subpackage/internal/systems/system1/domain"
	d2 "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-subpackage/internal/systems/system2/domain"
)

func ApiFn() {
	d1.DomainFn()
	d2.DomainFn()
}
