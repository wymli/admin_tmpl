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
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type ViewGenerator struct {
	SrcDir   string
	DestDir  string
	GinDir   string
	AxiosDir string
	ModPath  string
}

type App struct {
	Name     string
	FileName string
	Routes   []Route
	Handlers []Handler
	Models   []Model
}

type Route struct {
	Method   string
	Path     string
	FuncName string
}

type Handler struct {
	FileName string
	FuncName string
	Inputs   []FuncParam
	Outputs  []FuncParam
	Route    Route
}

type Model struct {
	Name string
}

func (g *ViewGenerator) Parse() (*App, error) {
	if err := checkDir(g.SrcDir); err != nil {
		return nil, err
	}

	modelStructMap := map[string]Struct{}
	app := &App{}
	modelPkgName := "model"

	os.ReadFile(path.Join(g.SrcDir, "model.go"))

	// 1. 解析所有请求/响应的model结构
	if err := filepath.Walk(g.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".gen.go") {
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
			case *ast.TypeSpec:
				switch x.Type.(type) {
				case *ast.StructType:
					modelStructMap[x.Name.Name] = Struct{
						FileName: filepath.Base(path),
						Name:     x.Name.Name,
					}
				}
			}
			return true
		})
		return nil
	}); err != nil {
		return nil, err
	}

	for k := range modelStructMap {
		app.Models = append(app.Models, Model{Name: k})
	}

	// 2. 解析 route/handlers
	if err := filepath.Walk(g.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".gen.go") {
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
								if _, ok := modelStructMap[paramType]; ok {
									pkgName = modelPkgName
									pkgPath = filepath.Join(g.ModPath, g.SrcDir)
								}
								fmt.Println("paramType:", paramType, modelStructMap[paramType], modelStructMap)

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
						handlerFile := ""

						route := Route{FuncName: methodName}
						for _, c := range method.Comment.List {
							items := strings.Split(strings.TrimPrefix(c.Text, "//"), " ")
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
									route.Method = strings.ToUpper(v)
								case "path":
									route.Path = v
								case "file":
									if !strings.HasSuffix(v, ".go") {
										v = v + ".go"
									}
									handlerFile = v
								}
							}
						}

						app.Routes = append(app.Routes, route)

						app.Handlers = append(app.Handlers, Handler{
							FileName: handlerFile,
							FuncName: methodName,
							Inputs:   inputs,
							Outputs:  outputs,
							Route:    route,
						})
					}
				}
			}
			return true
		})

		return nil
	}); err != nil {
		return nil, err
	}

	fmt.Println("app:", jsonIndent(app))

	return app, nil
}

func (g *ViewGenerator) Generate(app *App) error {
	fmt.Println("rendering")
	if err := g.renderModels(app); err != nil {
		return err
	}
	if err := g.renderHandlers(app.Handlers); err != nil {
		return err
	}
	if err := g.renderGinRoutes(app.Routes); err != nil {
		return err
	}
	if err := g.renderAxiosMethods(app.Handlers); err != nil {
		return err
	}
	fmt.Println("rendered")
	return nil
}

func (g *ViewGenerator) Format() error {
	return exec.Command("gofumpt", "-l", "-w", ".").Run()
}

func (g *ViewGenerator) renderModels(app *App) error {
	tmpl := NewTmplRenderOrDie(TmplRenderOpts{
		Tmpl: `
package model

type ViewModelI interface{
	StructName() string
} 

{{ .Models }}
`,
		RemovePrefixNewline: true,
	})

	content, err := tmpl.Render(struct {
		Models string
	}{
		Models: strings.Join(gsliceMap(app.Models, func(v Model) string {
			return fmt.Sprintf(`
func (v %s) StructName() string {
	return "%s"
}
`, v.Name, v.Name)
		}), "\n"),
	})
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(g.SrcDir, "model.gen.go"), []byte(content), 0o644)
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

// render to handlers/$resource.go
func (g *ViewGenerator) renderHandlers(handlers []Handler) error {
	for _, h := range handlers {
		fpath := filepath.Join(g.DestDir, h.FileName)

		_ = os.MkdirAll(filepath.Dir(fpath), 0o755)

		file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if len(bytes) == 0 {
			file.WriteString(fmt.Sprintf(`package handlers
import "%s/view/model"
`, g.ModPath))
		}

		if strings.Contains(string(bytes), fmt.Sprintf("func (app *App) %s(", h.FuncName)) {
			continue
		}

		inputStr := strings.Join(gsliceMap(h.Inputs, func(v FuncParam) string {
			return v.Name + " " + chooseIf(v.IsPointer, "*", "") + chooseIf(v.PkgName != "", v.PkgName+".", "") + v.Type
		}), ",")

		outputStr := strings.Join(gsliceMap(h.Outputs, func(v FuncParam) string {
			return v.Name + " " + chooseIf(v.IsPointer, "*", "") + chooseIf(v.PkgName != "", v.PkgName+".", "") + v.Type
		}), ",")

		d := h.Outputs[0]
		returnStr := fmt.Sprintf("return %s%s%s{}, nil", chooseIf(d.IsPointer, "&", ""), chooseIf(d.PkgName != "", d.PkgName+".", ""), d.Type)

		file.WriteString(fmt.Sprintf(`
func (app *App) %s(%s)(%s){
	%s
}
`, h.FuncName, inputStr, outputStr, returnStr))
	}
	return nil
}

// render to server/view_server/gin/route.gen.go
func (g *ViewGenerator) renderGinRoutes(routes []Route) error {
	tmpl := NewTmplRenderOrDie(TmplRenderOpts{
		Tmpl: `
package gin

import (
	"github.com/gin-gonic/gin"
	"{{ .ModPath }}/view/handlers"
)

func RegisterAppRouteGen(engine *gin.Engine, app *handlers.App) {
	{{ .Routes }}
}
`,
		RemovePrefixNewline: true,
	})

	content, err := tmpl.Render(struct {
		ModPath string
		Routes  string
	}{
		ModPath: g.ModPath,
		Routes: strings.Join(gsliceMap(routes, func(r Route) string {
			return fmt.Sprintf(`engine.Handle("%s", "%s", WrapGinHandler(handlers.WrapMiddlewares(app.%s)))`, r.Method, r.Path, r.FuncName)
		}), "\n\t"),
	})
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(g.GinDir, "route.gen.go"), []byte(content), 0o644)
}

func (g *ViewGenerator) renderAxiosMethods(handlers []Handler) error {
	createDirOrDie(g.AxiosDir)

	tmpl := NewTmplRenderOrDie(TmplRenderOpts{
		Tmpl: `
import { doRequest } from '@/utils/axios'

function {{ .FuncName }}({
    query = null,
    data = null,
    setLoadingFn = null,
    okFn = null,
    errFn = null,
}) {
    doRequest({
        method: "{{ .Method }}",
        url: "{{ .Path }}",
        query: query,
        data: data,
        setLoadingFn: setLoadingFn,
        okFn: okFn,
        errFn: errFn
    })
}
`,
		RemovePrefixNewline: true,
	})

	for _, v := range handlers {
		content, err := tmpl.Render(struct {
			FuncName string
			Method   string
			Path     string
		}{
			FuncName: v.Route.FuncName,
			Method:   v.Route.Method,
			Path:     v.Route.Path,
		})
		if err != nil {
			return err
		}

		if err := os.WriteFile(path.Join(g.AxiosDir, strings.ReplaceAll(v.FileName, ".go", ".tsx")), []byte(content), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func createDirOrDie(d string) {
	if err := os.MkdirAll(d, 0o755); err != nil {
		panic(err)
	}
}

func checkOrCreateFile(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
			return err
		}
	}

	if !fi.Mode().IsRegular() {
		return fmt.Errorf("path=%s not a reguler file", p)
	}

	return nil
}

func gsliceMap[T1 any, T2 any](v []T1, f func(T1) T2) []T2 {
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
