package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	goPackage := flag.String("package", "", "go package")
	declFile := flag.String("decl", "", "event declare go file (*.go)")
	genFile := flag.String("gen", "", "event generate go file (*.go)")
	corePackage := flag.String("core", "core", "core package")
	eventRegexp := flag.String("regexp", "^[eE]vent.+", "event regexp")
	exportEmit := flag.Bool("exportemit", true, "export emit")
	genassist := flag.String("genassist", "", "generate event assist code")

	flag.Parse()

	if *declFile == "" || filepath.Ext(*declFile) != ".go" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	if *goPackage == "" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	if *genFile == "" {
		*genFile = strings.TrimSuffix(*declFile, ".go") + "_gencode.go"
	} else if filepath.Ext(*genFile) != ".go" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	declFileData, err := ioutil.ReadFile(*declFile)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()

	fast, err := parser.ParseFile(fset, *declFile, declFileData, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	genCode := &bytes.Buffer{}

	fmt.Fprintf(genCode, `// Code generated by %s%s; DO NOT EDIT.
package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
		func() (args string) {
			for _, arg := range os.Args[1:] {
				args += " " + arg
			}
			return
		}(),
		*goPackage)

	importCode := &bytes.Buffer{}

	fmt.Fprintf(importCode, "\nimport (")

	if *corePackage != "" {
		fmt.Fprintf(importCode, `
	%s "github.com/pangdogs/core"`, *corePackage)
	}

	if *genassist != "" {
		fmt.Fprintf(importCode, `
	"github.com/pangdogs/core/container"`)
	}

	for _, imp := range fast.Imports {
		begin := fset.Position(imp.Pos())
		end := fset.Position(imp.End())

		impStr := string(declFileData[begin.Offset:end.Offset])

		if *corePackage != "" && strings.Contains(impStr, "github.com/pangdogs/core") {
			continue
		}

		if *genassist != "" && strings.Contains(impStr, "github.com/pangdogs/core/container") {
			continue
		}

		fmt.Fprintf(importCode, "\n\t%s", impStr)
	}

	fmt.Fprintf(importCode, "\n)\n")

	if importCode.Len() > 12 {
		fmt.Fprintf(genCode, importCode.String())
	}

	exp, err := regexp.Compile(*eventRegexp)
	if err != nil {
		panic(err)
	}

	exportEmitStr := "emit"

	if *exportEmit {
		exportEmitStr = "Emit"
	}

	type EventInfo struct {
		Name    string
		Comment string
	}

	var events []EventInfo

	ast.Inspect(fast, func(node ast.Node) bool {
		ts, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if !ts.Name.IsExported() {
			return true
		}

		eventName := ts.Name.Name
		var eventComment string

		for _, comment := range fast.Comments {
			if fset.Position(comment.End()).Line+1 == fset.Position(node.Pos()).Line {
				eventComment = comment.Text()
				break
			}
		}

		if !exp.MatchString(eventName) {
			return true
		}

		eventIFace, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		if eventIFace.Methods.NumFields() <= 0 {
			return true
		}

		eventFuncField := eventIFace.Methods.List[0]

		if len(eventFuncField.Names) <= 0 {
			return true
		}

		eventFuncName := eventFuncField.Names[0].Name

		eventFunc, ok := eventFuncField.Type.(*ast.FuncType)
		if !ok {
			return true
		}

		eventFuncParamsDecl := ""
		eventFuncParams := ""

		if eventFunc.Params != nil {
			for i, param := range eventFunc.Params.List {
				paramName := ""

				for _, pn := range param.Names {
					if paramName != "" {
						paramName += ", "
					}
					paramName += pn.Name
				}

				if paramName == "" {
					paramName = fmt.Sprintf("p%d", i)
				}

				if eventFuncParams != "" {
					eventFuncParams += ", "
				}
				eventFuncParams += paramName

				begin := fset.Position(param.Type.Pos())
				end := fset.Position(param.Type.End())

				eventFuncParamsDecl += fmt.Sprintf(", %s %s", paramName, declFileData[begin.Offset:end.Offset])
			}
		}

		eventFuncTypeParamsDecl := ""
		eventFuncTypeParams := ""

		if ts.TypeParams != nil {
			for i, typeParam := range ts.TypeParams.List {
				typeParamName := ""

				for _, pn := range typeParam.Names {
					if typeParamName != "" {
						typeParamName += ", "
					}
					typeParamName += pn.Name
				}

				if typeParamName == "" {
					typeParamName = fmt.Sprintf("p%d", i)
				}

				if eventFuncTypeParams != "" {
					eventFuncTypeParams += ", "
				}
				eventFuncTypeParams += typeParamName

				begin := fset.Position(typeParam.Type.Pos())
				end := fset.Position(typeParam.Type.End())

				if eventFuncTypeParamsDecl != "" {
					eventFuncTypeParamsDecl += ", "
				}
				eventFuncTypeParamsDecl += fmt.Sprintf("%s %s", typeParamName, declFileData[begin.Offset:end.Offset])
			}
		}

		if eventFuncTypeParamsDecl != "" {
			eventFuncTypeParamsDecl = fmt.Sprintf("[%s]", eventFuncTypeParamsDecl)
		}

		if eventFuncTypeParams != "" {
			eventFuncTypeParams = fmt.Sprintf("[%s]", eventFuncTypeParams)
		}

		_corePackage := ""
		if *corePackage != "" {
			_corePackage = *corePackage + "."
		}

		if eventFunc.Results.NumFields() > 0 {
			eventRet, ok := eventFunc.Results.List[0].Type.(*ast.Ident)
			if !ok {
				return true
			}

			if eventRet.Name != "bool" {
				return true
			}

			fmt.Fprintf(genCode, `
func %[8]s%[1]s%[6]s(event %[5]sIEvent%[3]s) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate %[5]sFastIFace) bool {
		return %[5]sFast2IFace[%[1]s%[7]s](delegate).%[2]s(%[4]s)
	})
}
`, eventName, eventFuncName, eventFuncParamsDecl, eventFuncParams, _corePackage, eventFuncTypeParamsDecl, eventFuncTypeParams, exportEmitStr)

		} else {

			fmt.Fprintf(genCode, `
func %[8]s%[1]s%[6]s(event %[5]sIEvent%[3]s) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate %[5]sFastIFace) bool {
		%[5]sFast2IFace[%[1]s%[7]s](delegate).%[2]s(%[4]s)
		return true
	})
}
`, eventName, eventFuncName, eventFuncParamsDecl, eventFuncParams, _corePackage, eventFuncTypeParamsDecl, eventFuncTypeParams, exportEmitStr)
		}

		events = append(events, EventInfo{
			Name:    eventName,
			Comment: eventComment,
		})

		fmt.Println(eventName)

		return true
	})

	if *genassist != "" {
		var eventsCode string
		var eventsRecursionCode string

		_corePackage := ""
		if *corePackage != "" {
			_corePackage = *corePackage + "."
		}

		for _, event := range events {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, _corePackage)
		}

		fmt.Fprintf(genCode, `
type %[1]sInterface interface {
%[2]s}
`, *genassist, eventsCode)

		for i, event := range events {
			eventRecursion := "EventRecursion_Disallow"

			if strings.Contains(event.Comment, "[EventRecursion_Allow]") {
				eventRecursion = "EventRecursion_Allow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Discard]") {
				eventRecursion = "EventRecursion_Discard"
			}

			eventsRecursionCode += fmt.Sprintf("\tassist.eventTab[%d].Init(autoRecover, reportError, %s%s, hookCache, gcCollector)\n", i, _corePackage, eventRecursion)
		}

		var eventsAccessCode string

		for i, event := range events {
			eventsAccessCode += fmt.Sprintf(`
func (assist *%s) %s() %sIEvent {
	return &assist.eventTab[%d]
}
`, *genassist, event.Name, _corePackage, i)
		}

		fmt.Fprintf(genCode, `
type %[1]s struct {
	eventTab [%[2]d]%[4]sEvent
}

func (assist *%[1]s) Init(autoRecover bool, reportError chan error, hookCache *container.Cache[%[4]sHook], gcCollector container.GCCollector) {
%[3]s}

func (assist *%[1]s) Shut() {
	for i := range assist.eventTab {
		assist.eventTab[i].Clear()
	}
}
%[5]s
`, *genassist, len(events), eventsRecursionCode, _corePackage, eventsAccessCode)
	}

	if err := ioutil.WriteFile(*genFile, genCode.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
