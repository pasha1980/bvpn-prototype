package tests

import (
	"bvpn-prototype/tests/testing"
)

const (
	TypeInfrastructure TestType = iota
	TypeFeature
	TypeUnit
)

type TestType int

func Run(types ...TestType) *testing.R {
	if len(types) == 0 {
		types = []TestType{
			TypeInfrastructure,
			TypeFeature,
			TypeUnit,
		}
	}

	prepareTestDI()
	masterResult := &testing.R{}

	for _, testType := range types {
		var list []func(t *testing.T) error
		switch testType {
		case TypeInfrastructure:
			list = infrastructureTestList
			break
		case TypeFeature:
			list = featureTestsList
			break
		case TypeUnit:
			list = unitTestsList
			break
		}

		result := do(list)
		masterResult.MergeResults(result)
	}

	return masterResult
}

func prepareTestDI() {
	// todo
}

func do(list []func(t *testing.T) error) *testing.R {
	result := &testing.R{}
	for _, f := range list {
		t := testing.T{}
		err := f(&t)
		if err != nil {
			result.Errors++
			continue
		}

		result.Add(&t)

	}

	return result
}
