package systems

import (
	api1 "github.com/quantumcycle/go-import-checks/validator/tests/allowed/allow-root-package-restrict-same-subpackage/internal/systems/system1/api"
	api2 "github.com/quantumcycle/go-import-checks/validator/tests/allowed/allow-root-package-restrict-same-subpackage/internal/systems/system2/api"
)

func SubsystemCall() {
	// This is allowed because the root systems package is allowed to import any all the subpackages apis
	api1.ApiFn()
	api2.ApiFn()
}
