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
	corePackage := flag.String("core", "core", "core package")
	eventRegexp := flag.String("regexp", "^[eE]vent.+", "event regexp")
	declFile := flag.String("decl", "", "event declare go file (*.go)")
	emitGOPackage := flag.String("emit_package", "", "emit event go package")
	emitGenFile := flag.String("gen_emit_dir", "", "generate emit event go file (*.go) dir")
	exportEmit := flag.Bool("export_emit", true, "export emit")
	genAssistCode := flag.String("gen_assist_code", "", "generate event assist code")
	assistGOPackage := flag.String("assist_package", "", "event assist go package")
	assistGenFile := flag.String("gen_assist_dir", "", "generate event assist go file (*.go) dir")

	flag.Parse()

	if *declFile == "" || filepath.Ext(*declFile) != ".go" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	if *emitGOPackage == "" {
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

	type EventInfo struct {
		Name    string
		Comment string
	}

	var events []EventInfo

	{
		if *emitGenFile == "" {
			*emitGenFile = strings.TrimSuffix(*declFile, ".go") + "_emit_code.go"
		} else {
			*emitGenFile = filepath.Dir(*declFile) + string(filepath.Separator) + filepath.Base(strings.TrimSuffix(*declFile, ".go")) + "_emit_code.go"
		}

		genEmitCodeBuff := &bytes.Buffer{}

		fmt.Fprintf(genEmitCodeBuff, `// Code generated by %s%s; DO NOT EDIT.
package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
			func() (args string) {
				for _, arg := range os.Args[1:] {
					args += " " + arg
				}
				return
			}(),
			*emitGOPackage)

		emitImportCode := &bytes.Buffer{}

		fmt.Fprintf(emitImportCode, "\nimport (")

		if *corePackage != "" {
			fmt.Fprintf(emitImportCode, `
	%s "github.com/pangdogs/core"`, *corePackage)
		}

		for _, imp := range fast.Imports {
			begin := fset.Position(imp.Pos())
			end := fset.Position(imp.End())

			impStr := string(declFileData[begin.Offset:end.Offset])

			if *corePackage != "" && strings.Contains(impStr, "github.com/pangdogs/core") {
				continue
			}

			fmt.Fprintf(emitImportCode, "\n\t%s", impStr)
		}

		fmt.Fprintf(emitImportCode, "\n)\n")

		if emitImportCode.Len() > 12 {
			fmt.Fprintf(genEmitCodeBuff, emitImportCode.String())
		}

		exp, err := regexp.Compile(*eventRegexp)
		if err != nil {
			panic(err)
		}

		exportEmitStr := "emit"

		if *exportEmit {
			exportEmitStr = "Emit"
		}

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

				fmt.Fprintf(genEmitCodeBuff, `
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

				fmt.Fprintf(genEmitCodeBuff, `
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

		if err := ioutil.WriteFile(*emitGenFile, genEmitCodeBuff.Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	}

	if *genAssistCode != "" {
		if *assistGOPackage == "" {
			*assistGOPackage = *emitGOPackage
		}

		if *assistGenFile == "" {
			*assistGenFile = strings.TrimSuffix(*declFile, ".go") + "_assist_code.go"
		} else {
			*assistGenFile = filepath.Dir(*declFile) + string(filepath.Separator) + filepath.Base(strings.TrimSuffix(*declFile, ".go")) + "_assist_code.go"
		}

		genAssistCodeBuff := &bytes.Buffer{}

		fmt.Fprintf(genAssistCodeBuff, `// Code generated by %s%s; DO NOT EDIT.
package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
			func() (args string) {
				for _, arg := range os.Args[1:] {
					args += " " + arg
				}
				return
			}(),
			*assistGOPackage)

		assistImportCode := &bytes.Buffer{}

		fmt.Fprintf(assistImportCode, "\nimport (")

		if *corePackage != "" {
			fmt.Fprintf(assistImportCode, `
	%s "github.com/pangdogs/core"`, *corePackage)
		}

		fmt.Fprintf(assistImportCode, `
	"github.com/pangdogs/core/container"`)

		for _, imp := range fast.Imports {
			begin := fset.Position(imp.Pos())
			end := fset.Position(imp.End())

			impStr := string(declFileData[begin.Offset:end.Offset])

			if *corePackage != "" && strings.Contains(impStr, "github.com/pangdogs/core") {
				continue
			}

			if *genAssistCode != "" && strings.Contains(impStr, "github.com/pangdogs/core/container") {
				continue
			}

			fmt.Fprintf(assistImportCode, "\n\t%s", impStr)
		}

		fmt.Fprintf(assistImportCode, "\n)\n")

		if assistImportCode.Len() > 12 {
			fmt.Fprintf(genAssistCodeBuff, assistImportCode.String())
		}

		var eventsCode string
		var eventsRecursionCode string

		_corePackage := ""
		if *corePackage != "" {
			_corePackage = *corePackage + "."
		}

		for _, event := range events {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, _corePackage)
		}

		fmt.Fprintf(genAssistCodeBuff, `
type I%[1]s interface {
%[2]s}
`, *genAssistCode, eventsCode)

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
const %[2]sID int = %[4]d

func (assist *%[1]s) %[2]s() %[3]sIEvent {
	return &assist.eventTab[%[2]sID]
}
`, *genAssistCode, event.Name, _corePackage, i)
		}

		fmt.Fprintf(genAssistCodeBuff, `
type %[1]s struct {
	eventTab [%[2]d]%[4]sEvent
}

func (assist *%[1]s) EventTab(id int) %[4]sIEvent {
	return &assist.eventTab[id]
}

func (assist *%[1]s) Init(autoRecover bool, reportError chan error, hookCache *container.Cache[%[4]sHook], gcCollector container.GCCollector) {
%[3]s}

func (assist *%[1]s) Open() {
	for i := range assist.eventTab {
		assist.eventTab[i].Open()
	}
}

func (assist *%[1]s) Close() {
	for i := range assist.eventTab {
		assist.eventTab[i].Close()
	}
}

func (assist *%[1]s) Clear() {
	for i := range assist.eventTab {
		assist.eventTab[i].Clear()
	}
}
%[5]s
`, *genAssistCode, len(events), eventsRecursionCode, _corePackage, eventsAccessCode)

		if err := ioutil.WriteFile(*assistGenFile, genAssistCodeBuff.Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
