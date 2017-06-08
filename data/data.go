// Package data contains implementations that provide and consume NationStates
// data via an intermediate tree structure.
package data

import (
	"github.com/wcerfgba/nationstates-xlsx/util"
)

// Provider must provide an Result on request, and are Configurable to make Get
// nullary.
type Provider interface {

	// Get returns a Result.
	Get() Result

	// See util.Configurable.
	Configure(util.Configuration)
}

// Result is the parent node in a tree of StringTreeNodes.
type Result util.StringTreeNode

// Outputter must accept a Result on request, and are Configurable to make
// Output unary.
type Outputter interface {
	// Output does something with the Result, and may error.
	Output(Result) error

	// See util.Configurable.
	Configure(util.Configuration)
}
