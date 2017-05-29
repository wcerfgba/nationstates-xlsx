// Package data contains implementations that provide and conume NationStates
// data via an intermediate tree structure.
package data

import (
	"github.com/wcerfgba/nationstates-xlsx/util"
)

// Provider must provide an Result on request, and are configured
// beforehand to make Get niladic.
type Provider interface {
	Get() Result
	Configure(util.Configuration)
}

// Result is the parent node in a tree of StringTreeNodes.
type Result util.StringTreeNode

func NewResult() Result {
	return &util.StringTreeNode20170521{}
}

// Outputter must accept a Result on request, and are configured
// beforehand to make Output uniadic.
type Outputter interface {
	Output(Result) error
	Configure(util.Configuration)
}
