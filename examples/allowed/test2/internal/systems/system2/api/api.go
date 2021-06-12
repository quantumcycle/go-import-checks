package api

import (
	d1 "github.com/matdurand/project1/internal/systems/system1/domain"
	d2 "github.com/matdurand/project1/internal/systems/system2/domain"
)

func ApiFn() {
	d1.DomainFn()
	d2.DomainFn()
}
