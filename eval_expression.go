package main

import (
	"github.com/Knetic/govaluate"
)

// EvalIfExpression - check whether the specified expression evaluates as true
func EvalIfExpression(expr string, params map[string]interface{}) bool {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err == nil {
		result, err := expression.Evaluate(params)
		if err == nil {
			return result.(bool) == true
		}
	}
	return false
}

// GetExpressionValue - compile the specified expression
func GetExpressionValue(expr string, params map[string]interface{}) string {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err == nil {
		result, err := expression.Evaluate(params)
		if err == nil {
			return result.(string)
		}
	}
	return expr
}
