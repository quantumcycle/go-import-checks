package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readYamlCfg(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func assertValidation(t *testing.T, path string, expectedErrs []ValidationError) {
	cfg, err := readYamlCfg(filepath.Join(path, "config.yaml"))
	assert.Nil(t, err)

	validationErrs, err := Validate(path, cfg.Checks, false)

	assert.Nil(t, err)
	assert.Len(t, validationErrs, len(expectedErrs))

	for _, e := range validationErrs {
		found := false
		for _, expErr := range expectedErrs {
			if (e.Path == expErr.Path && e.Reason == expErr.Reason &&  e.ImportPath == expErr.ImportPath) {
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

func validationErrsPrettyPrint(errs []ValidationError) string {
	strs := []string{}
	for _, e := range errs {
		strs = append(strs, validationErrPrettyPrint(e))
	}
	return strings.Join(strs, "\n")
}

func validationErrPrettyPrint(e ValidationError) string {
	if (e.Reason == ReasonRejected) {
		return fmt.Sprintf("[Import {%s} explicitly rejected in {%s}]", e.ImportPath, e.Path)
	} else if (e.Reason == ReasonNotAllow) {
		return fmt.Sprintf("[Import {%s} not allowed in {%s}]", e.ImportPath, e.Path)
	} else {
		panic("Unexpected reason")
	}
}

func TestAllowSamePackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-same-package", []ValidationError{
		{
			Path: "pkg/p1/p1.go",
			Reason: ReasonNotAllow,
			ImportPath: "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-package/internal",
		},
	})
}

func TestAllowSameSubPackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-same-subpackage", []ValidationError{
		{
			Path: "internal/systems/system2/api/api.go",
			Reason: ReasonNotAllow,
			ImportPath: "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-same-subpackage/internal/systems/system1/domain",
		},
	})
}

func TestAllowWildcardSubpackage(t *testing.T) {
	assertValidation(t, "tests/allowed/restrict-wildcard-subpackage", []ValidationError{
		{
			Path: "internal/server.go",
			Reason: ReasonNotAllow,
			ImportPath: "github.com/matdurand/go-import-checks/validator/tests/allowed/restrict-wildcard-subpackage/components/component1/services",
		},
	})
}

func TestRejectAnotherPackage(t *testing.T) {
	assertValidation(t, "tests/reject/reject-another-package", []ValidationError{
		{
			Path:       "pkg/p1/p1.go",
			Reason:     ReasonRejected,
			ImportPath: "github.com/matdurand/go-import-checks/validator/tests/reject/reject-another-package/internal",
		},
	})
}

