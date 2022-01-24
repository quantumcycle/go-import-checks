package internal

import (
	"github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-wildcard-subpackage/components/component1/services"
	"github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-wildcard-subpackage/components/component2/api"
)

func InitServer() {
	services.ServiceMethod()
	api.ApiMethod()
}