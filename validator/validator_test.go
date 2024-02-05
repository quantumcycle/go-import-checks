package validator_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/quantumcycle/go-import-checks/validator"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func readYamlCfg(path string) (validator.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return validator.Config{}, err
	}

	cfg := validator.Config{}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return validator.Config{}, err
	}

	return cfg, nil
}

func assertValidation(t *testing.T, path string, expectedErrs []validator.ValidationError) {
	cfg, err := readYamlCfg(filepath.Join(path, "config.yaml"))
	assert.Nil(t, err)

	validationErrs, err := validator.Validate(path, cfg.Checks, false)

	assert.Nil(t, err)
	assert.Len(t, validationErrs, len(expectedErrs))

	for _, e := range validationErrs {
		found := false
		for _, expErr := range expectedErrs {
			if e.Path == expErr.Path && e.Reason == expErr.Reason && e.ImportPath == expErr.ImportPath {
				found = true
			}
		}
		if !found {
			assert.Fail(t, fmt.Sprintf("Cannot find validation error: %s \nin {\n%s\n}",
				validationErrPrettyPrint(e),
				validationErrsPrettyPrint(expectedErrs)))
		}
	}
}

func validationErrsPrettyPrint(errs []validator.ValidationError) string {
	strs := []string{}
	for _, e := range errs {
		strs = append(strs, validationErrPrettyPrint(e))
	}
	return strings.Join(strs, "\n")
}

func validationErrPrettyPrint(e validator.ValidationError) string {
	if e.Reason == validator.ReasonRejected {
		return fmt.Sprintf("[Import {%s} explicitly rejected in {%s}]", e.ImportPath, e.Path)
	} else if e.Reason == validator.ReasonNotAllow {
		return fmt.Sprintf("[Import {%s} not allowed in {%s}]", e.ImportPath, e.Path)
	} else {
		panic("Unexpected reason")
	}
}

func TestAllowSamePackageWithExclude(t *testing.T) {
	//p1.go has some error, but is excluded from the checks, so no errors should be reported
	assertValidation(t, "tests/allowed/restrict-but-exclude", []validator.ValidationError{})
}

func TestAllowSamePackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-same-package", []validator.ValidationError{
		{
			Path:       "pkg/p1/p1.go",
			Reason:     validator.ReasonNotAllow,
			ImportPath: "github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-same-package/internal",
		},
	})
}

func TestAllowSameSubPackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-same-subpackage", []validator.ValidationError{
		{
			Path:       "internal/systems/system2/api/api.go",
			Reason:     validator.ReasonNotAllow,
			ImportPath: "github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-same-subpackage/internal/systems/system1/domain",
		},
	})
}

func TestAllowRootCallSubpackages(t *testing.T) {
	assertValidation(t, "tests/allowed/allow-root-package-restrict-same-subpackage", []validator.ValidationError{
		{
			Path:       "internal/systems/system2/api/api.go",
			Reason:     validator.ReasonNotAllow,
			ImportPath: "github.com/quantumcycle/go-import-checks/validator/tests/allowed/allow-root-package-restrict-same-subpackage/internal/systems/system1/domain",
		},
		//no errors on init.go in the systems package. we want that since it should not be covered by the existing rule
	})
}

func TestAllowWildcardSubpackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-wildcard-subpackage", []validator.ValidationError{
		{
			Path:       "internal/server.go",
			Reason:     validator.ReasonNotAllow,
			ImportPath: "github.com/quantumcycle/go-import-checks/validator/tests/allowed/restrict-wildcard-subpackage/components/component1/services",
		},
	})
}

func TestRejectAnotherPackage(t *testing.T) {
	assertValidation(t, "tests/reject/reject-another-package", []validator.ValidationError{
		{
			Path:       "pkg/p1/p1.go",
			Reason:     validator.ReasonRejected,
			ImportPath: "github.com/quantumcycle/go-import-checks/validator/tests/reject/reject-another-package/internal",
		},
	})
}
