package validator

import (
	"fmt"
	"github.com/matdurand/go-import-checks/glob"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Rule struct {
	Prefix string   `yaml:"prefix"`
	Reject []string `yaml:"reject"`
	Allow  []string `yaml:"allow"`
}

type Check struct {
	Folder      string `yaml:"folder"`
	Subpackages bool   `yaml:"subpackages"`
	Rules       []Rule `yaml:"rules"`
}

func (chk Check) isApplicable(path string) bool {
	pathParts := strings.Split(path, "/")
	chkParts := strings.Split(chk.Folder, "/")
	for i, chkItem := range chkParts {
		if i >= len(pathParts) {
			return false
		}

		if strings.HasPrefix(chkItem, "$") || chkItem == "*" {
			continue
		}

		pathItem := pathParts[i]
		if chkItem != pathItem {
			return false
		}
	}

	return len(pathParts) == len(chkParts) || chk.Subpackages
}

func (chk Check) extractPathVariables(path string) map[string]string {
	vars := make(map[string]string, 0)
	pathParts := strings.Split(path, "/")
	chkParts := strings.Split(chk.Folder, "/")
	for i, chkItem := range chkParts {
		if strings.HasPrefix(chkItem, "$") {
			varName := chkItem[1:]
			vars[varName] = pathParts[i]
		}
	}
	return vars
}

func replaceVariables(pck string, variables map[string]string) string {
	if strings.Index(pck, "$") == -1 {
		return pck
	}

	parts := strings.Split(pck, "/")
	resolvedParts := make([]string, 0, len(parts))
	for _, p := range parts {
		var varName string
		if strings.HasPrefix(p, "!$") {
			varName = p[2:]
		}
		if strings.HasPrefix(p, "$") {
			varName = p[1:]
		}
		if varName != "" {
			varValue := variables[varName]
			if varValue == "" {
				panic(
					"Some rules are using a variable named [" + varName + "] which is not defined in the matchinng folder expression",
				)
			}
			resolvedParts = append(resolvedParts, strings.Replace(p, "$"+varName, varValue, 1))
		} else {
			resolvedParts = append(resolvedParts, p)
		}
	}
	return strings.Join(resolvedParts, "/")
}


func isPackageMatchingExpression(pck string, pckExpression string) bool {
	g, err := glob.NewGlob(pckExpression)
	if (err != nil) {
		panic(err)
	}

	return g.Match(pck)
}

func validateImport(
	path string,
	chk Check,
	r Rule,
	imprt string,
	variables map[string]string,
	debug bool,
) *ValidationError {
	if strings.HasPrefix(imprt, r.Prefix) {
		if len(r.Allow) > 0 {
			var allowed bool
			for _, allowSuffix := range r.Allow {
				if debug {
					fmt.Printf("DEBUG: \t\t-checking allowed rule [%s] on import [%s]:", allowSuffix, imprt)
				}
				pck := replaceVariables(r.Prefix+allowSuffix, variables)
				impAllowed := isPackageMatchingExpression(imprt, pck)
				if debug {
					if impAllowed {
						fmt.Printf(" allowed\n")
					} else {
						fmt.Printf(" not allowed\n")
					}
				}
				if impAllowed {
					allowed = true
					break
				}
			}
			if !allowed {
				return &ValidationError{
					Path:       path,
					Check:      chk,
					Rule:       r,
					ImportPath: imprt,
					Reason:     ReasonNotAllow,
				}
			}
		}

		if len(r.Reject) > 0 {
			for _, rejectSuffix := range r.Reject {
				if debug {
					fmt.Printf("DEBUG: \t\t-checking rejection rule [%s] on import [%s]:", rejectSuffix, imprt)
				}
				pck := replaceVariables(r.Prefix+rejectSuffix, variables)
				matching := isPackageMatchingExpression(imprt, pck)
				if debug {
					if matching {
						fmt.Printf(" rejected\n")
					} else {
						fmt.Printf(" not rejected\n")
					}
				}
				if matching {
					return &ValidationError{
						Path:       path,
						Check:      chk,
						Rule:       r,
						ImportPath: imprt,
						Reason:     ReasonRejected,
					}
				}
			}
		}
	}

	return nil
}

func (chk Check) ValidateImports(path string, imports []*ast.ImportSpec, debug bool) []ValidationError {
	errors := []ValidationError{}
	if debug {
		if chk.isApplicable(path) {
			fmt.Printf("DEBUG: \t\tcheck [%s] is applicable\n", chk.Folder)
		} else {
			fmt.Printf("DEBUG: \t\tcheck [%s] is not applicable, skipping\n", chk.Folder)
		}
	}
	if chk.isApplicable(path) {
		variables := chk.extractPathVariables(path)
		for _, r := range chk.Rules {
			for _, imp := range imports {
				impPath := imp.Path.Value
				impPath = strings.TrimSuffix(impPath, `"`)
				impPath = strings.TrimPrefix(impPath, `"`)
				err := validateImport(path, chk, r, impPath, variables, debug)
				if err != nil {
					errors = append(errors, *err)
				}
			}
		}
	}
	return errors
}

type Config struct {
	Checks []Check `yaml:"imports-checks"`
}

type Reason int

const (
	ReasonNotAllow = iota + 1
	ReasonRejected
)

type ValidationError struct {
	Path       string
	Check      Check
	Rule       Rule
	ImportPath string
	Reason     Reason
}

func validateChecks(checks []Check) {
	for _, c := range checks {
		if c.Folder == "" {
			panic("Some checks have an empty folder attribute")
		}

		for _, r := range c.Rules {
			if len(r.Allow) > 0 && len(r.Reject) > 0 {
				panic(
					"Rules for folder " + c.Folder + " have both reject and allow definition. Use either one but not both.",
				)
			}
		}
	}
}

func Validate(root string, checks []Check, debug bool) ([]ValidationError, error) {
	validateChecks(checks)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, err
	}

	errors := []ValidationError{}
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			panic(err)
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			panic(err)
		}

		ferrors := validateImports(relPath, f.Imports, checks, debug)
		errors = append(errors, ferrors...)

		return nil
	})

	return errors, err
}

func validateImports(path string, imp []*ast.ImportSpec, checks []Check, debug bool) []ValidationError {
	if debug {
		fmt.Printf("DEBUG: Validating imports for file [%s]\n", path)
	}
	errors := []ValidationError{}
	for _, chk := range checks {
		errs := chk.ValidateImports(path, imp, debug)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	return errors
}
