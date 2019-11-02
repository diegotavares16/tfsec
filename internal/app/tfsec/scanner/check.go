package scanner

import (
	"fmt"

	"github.com/liamg/tfsec/internal/app/tfsec/parser"
)

// CheckCode is a unique identifier for a check
type CheckCode string

// Check is a targeted security test which can be applied to terraform templates. It includes the types to run on e.g.
// "resource", and the labels to run on e.g. "aws_s3_bucket".
type Check struct {
	Code           CheckCode
	RequiredTypes  []string
	RequiredLabels []string
	CheckFunc      func(*Check, *parser.Block) []Result
}

// Run runs the check against the provided HCL block, including the hclEvalContext to evaluate expressions if it is
// provided.
func (check *Check) Run(block *parser.Block) []Result {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("WARNING: fatal error running check: %s\n", err)
		}
	}()
	return check.CheckFunc(check, block)
}

// IsRequiredForBlock returns true if the Check should be applied to the given HCL block
func (check *Check) IsRequiredForBlock(block *parser.Block) bool {

	if check.CheckFunc == nil {
		return false
	}

	if len(check.RequiredTypes) > 0 {
		var found bool
		for _, requiredType := range check.RequiredTypes {
			if block.Type() == requiredType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(check.RequiredLabels) > 0 {
		var found bool
		for _, requiredLabel := range check.RequiredLabels {
			if len(block.Labels()) > 0 && block.Labels()[0] == requiredLabel {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// NewResult creates a new Result, containing the given description and range
func (check *Check) NewResult(description string, r parser.Range) Result {
	return Result{
		Code:        check.Code,
		Description: description,
		Range:       r,
	}
}

func (check *Check) NewResultWithAnnotation(description string, r parser.Range, annotation string) Result {
	return Result{
		Code:            check.Code,
		Description:     description,
		Range:           r,
		RangeAnnotation: annotation,
	}
}
