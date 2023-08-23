package api

import (
	d1 "github.com/quantumcycle/go-import-checks/validator/tests/allowed/allow-root-package-restrict-same-subpackage/internal/systems/system1/domain"
	d2 "github.com/quantumcycle/go-import-checks/validator/tests/allowed/allow-root-package-restrict-same-subpackage/internal/systems/system2/domain"
)

func ApiFn() {
	d1.DomainFn()
	d2.DomainFn()
}
