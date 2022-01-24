package api

import "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-subpackage/internal/systems/system1/domain"

func ApiFn() {
	domain.DomainFn()
}
