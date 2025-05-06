package main

func main() {
	g := &ViewGenerator{
		SrcDir:   "./view/model",
		DestDir:  "./view/handlers",
		GinDir:   "./view/server/gin",
		AxiosDir: "./fe/src/dal",
		ModPath:  "github.com/wymli/xserver",
	}

	app, err := g.Parse()
	if err != nil {
		panic(err)
	}

	if err := g.Generate(app); err != nil {
		panic(err)
	}

	if err := g.Format(); err != nil {
		panic(err)
	}
}
