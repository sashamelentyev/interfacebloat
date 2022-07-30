package analyzer

import (
	"flag"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const InterfaceLenFlag = "interface-len"

const defaultInterfaceLen = 10

// New returns new interfacebloat analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "interfacebloat",
		Doc:      "A linter that checks length of interface.",
		Run:      run,
		Flags:    flags(),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.Int(InterfaceLenFlag, 10, "length of interface")
	return *flags
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	filter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	insp.Preorder(filter, func(node ast.Node) {
		i, ok := node.(*ast.InterfaceType)
		if !ok {
			return
		}
		interfaceLen := interfaceLen(pass, InterfaceLenFlag)
		if len(i.Methods.List) > interfaceLen {
			report(pass, node.Pos(), interfaceLen)
		}
	})

	return nil, nil
}

func interfaceLen(pass *analysis.Pass, name string) (interfaceLen int) {
	interfaceLen, ok := pass.Analyzer.Flags.Lookup(name).Value.(flag.Getter).Get().(int)
	if !ok {
		interfaceLen = defaultInterfaceLen
	}
	return
}

func report(pass *analysis.Pass, pos token.Pos, interfaceLen int) {
	pass.Reportf(pos, `length of interface greater than %d`, interfaceLen)
}
