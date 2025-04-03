package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ViewGenerator struct {
	SrcDir  string
	DestDir string
	ModPath string
}

func (g *ViewGenerator) Generate() error {
	if err := checkDir(g.SrcDir); err != nil {
		return err
	}

	handlers := []Handler{}
	structMap := map[string]Struct{}
	app := &App{}
	modelPkgName := ""

	if err := filepath.Walk(g.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Printf("parsing file=%s\n", path)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			panic(err)
		}

		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				modelPkgName = x.Name.Name
			case *ast.TypeSpec:
				switch x.Type.(type) {
				case *ast.StructType:
					structMap[x.Name.Name] = Struct{
						FileName: filepath.Base(path),
						Name:     x.Name.Name,
					}
				}
			}
			return true
		})
		return nil
	}); err != nil {
		return err
	}

	if err := filepath.Walk(g.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Printf("parsing file=%s\n", path)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		ast.Print(fset, f)

		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				switch z := x.Type.(type) {
				case *ast.InterfaceType:
					app.Name = x.Name.Name
					app.FileName = filepath.Base(path)

					fmt.Printf("%s:\tfind app interface: %s in %s\n", fset.Position(n.Pos()), app.Name, app.FileName)
					for _, method := range z.Methods.List {
						methodType, ok := method.Type.(*ast.FuncType)
						if !ok {
							continue
						}
						if len(method.Names) != 1 {
							panic("one name, one type")
						}
						methodName := method.Names[0].Name

						parseParams := func(fields []*ast.Field) []FuncParam {
							res := []FuncParam{}
							for _, param := range fields {
								if len(param.Names) != 1 {
									panic("one name, one type")
								}

								paramName := param.Names[0].Name
								paramType := "<unknown>"
								paramIsPointer := false

								switch pt := param.Type.(type) {
								case *ast.StarExpr:
									paramType = pt.X.(*ast.Ident).Name
									paramIsPointer = true
								case *ast.Ident:
									paramType = pt.Name
								}

								pkgName := ""
								pkgPath := ""
								if _, ok := structMap[paramType]; ok {
									pkgName = modelPkgName
									pkgPath = filepath.Join(g.ModPath, g.SrcDir)
								}
								fmt.Println("paramType:", paramType, structMap[paramType], structMap)

								res = append(res, FuncParam{
									Name:      paramName,
									Type:      paramType,
									IsPointer: paramIsPointer,
									PkgName:   pkgName,
									PkgPath:   pkgPath,
								})
							}
							return res
						}

						inputs := parseParams(methodType.Params.List)
						outputs := parseParams(methodType.Results.List)
						handlerFile := decideHandlerFile(structMap, inputs, outputs)

						handlers = append(handlers, Handler{
							PackageName: "",
							FileName:    handlerFile,
							FuncName:    methodName,
							Inputs:      inputs,
							Outputs:     outputs,
						})

						route := Route{}
						for _, c := range method.Comment.List {
							items := strings.Split(strings.TrimLeft(c.Text, "//"), " ")
							for _, item := range items {
								if item == "" {
									continue
								}
								k, v, ok := strings.Cut(item, "=")
								if !ok {
									continue
								}
								switch k {
								case "method":
									route.Method = v
								case "path":
									route.Path = v
								}
							}
						}

						app.Routes = append(app.Routes, route)
					}
				}
			}
			return true
		})

		return nil
	}); err != nil {
		return err
	}

	fmt.Println("app:", jsonIndent(app))
	fmt.Println("rendering")
	if err := g.renderApp(app); err != nil {
		panic(err)
	}
	if err := g.renderHandlers(handlers); err != nil {
		panic(err)
	}
	fmt.Println("rendered")

	return nil
}

func decideHandlerFile(structMap map[string]Struct, inputs []FuncParam, outputs []FuncParam) string {
	params := []FuncParam{}
	params = append(params, inputs...)
	params = append(params, outputs...)

	fileCountMap := map[string]int{}
	for _, p := range params {
		if s, ok := structMap[p.Type]; ok {
			fileCountMap[s.FileName] += 1
		}
	}

	maxValue := 0
	maxKey := ""
	for k, v := range fileCountMap {
		if v >= maxValue {
			maxValue = v
			maxKey = k
		}
	}

	return maxKey
}

func checkDir(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("srcdir=%s is not a dir", p)
	}

	return nil
}

type App struct {
	Name     string
	FileName string
	Routes   []Route
}

type Route struct {
	Method string
	Path   string
}

type Handler struct {
	PackageName string
	FileName    string
	FuncName    string
	Inputs      []FuncParam
	Outputs     []FuncParam
}

func (h *Handler) GetImportPkgs() []string {
	params := []FuncParam{}
	params = append(params, h.Inputs...)
	params = append(params, h.Outputs...)

	res := map[string]struct{}{}
	for _, p := range params {
		if p.PkgPath == "" {
			continue
		}
		res[p.PkgPath] = struct{}{}
	}

	return gmapKeys(res)
}

type FuncParam struct {
	Name      string
	Type      string
	IsPointer bool
	PkgName   string
	PkgPath   string
}

func jsonIndent(v any) string {
	bytes, _ := json.MarshalIndent(v, "", "\t")
	return string(bytes)
}

type Struct struct {
	FileName string
	Name     string
}

func (g *ViewGenerator) renderApp(app *App) error {
	fpath := filepath.Join(g.DestDir, app.FileName)
	_ = os.MkdirAll(filepath.Dir(fpath), 0755)
	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		file.WriteString("package " + filepath.Base(g.DestDir) + "\n")
		file.WriteString("type App struct{}" + "\n")
	}
	return nil
}

func (g *ViewGenerator) renderHandlers(handlers []Handler) error {
	for _, h := range handlers {
		fpath := filepath.Join(g.DestDir, h.FileName)

		_ = os.MkdirAll(filepath.Dir(fpath), 0755)

		file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if strings.Contains(string(bytes), fmt.Sprintf("func (app *App) %s(", h.FuncName)) {
			continue
		}

		importPkgs := h.GetImportPkgs()
		fmt.Printf("%#v\n", importPkgs)
		importStr := strings.Join(gsliceTo(importPkgs, func(v string) string { return fmt.Sprintf("import \"%s\"", v) }), "\n")

		inputStr := strings.Join(gsliceTo(h.Inputs, func(v FuncParam) string {
			return v.Name + " " + chooseIf(v.IsPointer, "*", "") + chooseIf(v.PkgName != "", v.PkgName+".", "") + v.Type
		}), ",")

		outputStr := strings.Join(gsliceTo(h.Outputs, func(v FuncParam) string {
			return v.Name + " " + chooseIf(v.IsPointer, "*", "") + chooseIf(v.PkgName != "", v.PkgName+".", "") + v.Type
		}), ",")

		d := h.Outputs[0]
		returnStr := fmt.Sprintf("return %s%s%s{}, nil", chooseIf(d.IsPointer, "&", ""), chooseIf(d.PkgName != "", d.PkgName+".", ""), d.Type)

		if len(bytes) == 0 {
			file.WriteString("package " + filepath.Base(g.DestDir))
		}

		file.WriteString(fmt.Sprintf(`
%s
func (app *App) %s(%s)(%s){
	%s
}
`, importStr, h.FuncName, inputStr, outputStr, returnStr))
	}
	return nil
}

func checkOrCreateFile(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		if err := os.WriteFile(p, []byte{}, 0x644); err != nil {
			return err
		}
	}

	if !fi.Mode().IsRegular() {
		return fmt.Errorf("path=%s not a reguler file", p)
	}

	return nil
}

func gsliceTo[T1 any, T2 any](v []T1, f func(T1) T2) []T2 {
	res := []T2{}
	for _, vv := range v {
		res = append(res, f(vv))
	}
	return res
}

func chooseIf[T any](cond bool, onTrue T, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

func gmapKeys[K comparable, V any](v map[K]V) []K {
	res := []K{}
	for k := range v {
		res = append(res, k)
	}
	return res
}
