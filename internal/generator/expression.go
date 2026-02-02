package generator

import (
	"fmt"
	"math"
	"math/rand"
)

// Expression interface for our AST
type Expression interface {
	Eval(x, y, w, h float64) float64
	String() string
}

// Constant value node
type ValNode struct {
	Value float64
}

func (n ValNode) Eval(x, y, w, h float64) float64 { return n.Value }
func (n ValNode) String() string               { return fmt.Sprintf("%.2f", n.Value) }

// Variable node (X or Y coordinate, normalized 0-1)
type VarNode struct {
	Name string // "x" or "y"
}

func (n VarNode) Eval(x, y, w, h float64) float64 {
	if n.Name == "x" {
		return x / w
	}
	return y / h
}
func (n VarNode) String() string { return n.Name }

// Binary Operation node
type OpNode struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n OpNode) Eval(x, y, w, h float64) float64 {
	l := n.Left.Eval(x, y, w, h)
	r := n.Right.Eval(x, y, w, h)
	switch n.Op {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		if r == 0 { return 0 }
		return l / r
	case "%":
		return math.Mod(math.Abs(l), math.Abs(r) + 0.001)
	case "xor":
		// simulate bitwise xor on floats by scaling up
		return float64(int(l*255) ^ int(r*255)) / 255.0
	}
	return 0
}

func (n OpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op, n.Right.String())
}

// Unary Operation node (Sin, Cos, Abs)
type UnaryNode struct {
	Op   string
	Expr Expression
}

func (n UnaryNode) Eval(x, y, w, h float64) float64 {
	v := n.Expr.Eval(x, y, w, h)
	switch n.Op {
	case "sin":
		return math.Sin(v * math.Pi * 2)
	case "cos":
		return math.Cos(v * math.Pi * 2)
	case "abs":
		return math.Abs(v)
	case "tan":
		return math.Tan(v * math.Pi * 2)
	}
	return v
}

func (n UnaryNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Op, n.Expr.String())
}

// GenerateRandomExpression builds a random AST
func GenerateRandomExpression(depth int) Expression {
	if depth <= 0 || (depth > 1 && rand.Float64() < 0.2) {
		// Terminal node
		if rand.Float64() < 0.5 {
			return ValNode{Value: rand.Float64() * 5}
		}
		if rand.Float64() < 0.5 {
			return VarNode{Name: "x"}
		}
		return VarNode{Name: "y"}
	}

	// Operator node
	r := rand.Float64()
	if r < 0.6 {
		// Binary
		ops := []string{"+", "-", "*", "/", "%", "xor"}
		op := ops[rand.Intn(len(ops))]
		return OpNode{
			Op:    op,
			Left:  GenerateRandomExpression(depth - 1),
			Right: GenerateRandomExpression(depth - 1),
		}
	} else {
		// Unary
		ops := []string{"sin", "cos", "abs", "tan"}
		op := ops[rand.Intn(len(ops))]
		return UnaryNode{
			Op:   op,
			Expr: GenerateRandomExpression(depth - 1),
		}
	}
}
