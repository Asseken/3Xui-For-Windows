# 3X-UI front-end separation To Windows


## Reproduce from 3x ui
## Temporarily migrate and store front-end files in the etc/html directory，The file does not contain front-end files.
## ! modify 

## modify


--------------------------------------------------------
func (s *Server) getHtmlFiles() ([]string, error) {
	files := make([]string, 0)
	//dir, _ := os.Getwd() //3x-ui
	err := filepath.WalkDir(s.webDir, func(path string, d fs.DirEntry, err error) error { //new
		//err := fs.WalkDir(os.DirFS(dir), "web/html", func(path string, d fs.DirEntry, err error) error { //3x-ui
       //
}
----------------------------------------------------------
func (s *Server) getHtmlTemplate(funcMap template.FuncMap) (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	err := filepath.WalkDir(s.webDir, func(path string, d fs.DirEntry, err error) error { //new
		//err := fs.WalkDir(htmlFS, "html", func(path string, d fs.DirEntry, err error) error { //3x-ui
		if err != nil {
			return err
		}

		if d.IsDir() {
			newT, err := t.ParseGlob(filepath.Join(path, "*.html")) //new
			//newT, err := t.ParseFS(htmlFS, path+"/*.html") //3x-ui
			if err != nil {
				// ignore
				return nil
			}
			//
}
----------------------------------------------------------
func (s *Server) initRouter() (*gin.Engine, error) {
// set static files and template !!!thins modify
	if config.IsDebug() {
		// for development
		files, err := s.getHtmlFiles()
		if err != nil {
			return nil, err
		}
		engine.LoadHTMLFiles(files...)
		engine.Static("/assets", filepath.Join(s.webDir, "assets")) //new
		//engine.StaticFS(basePath+"assets", http.FS(os.DirFS("web/assets"))) //3x-ui。old
	} else {
		// for production
		template, err := s.getHtmlTemplate(engine.FuncMap)
		if err != nil {
			return nil, err
		}
		engine.SetHTMLTemplate(template)
		engine.Static("/assets", filepath.Join(s.webDir, "assets")) //new
		//engine.StaticFS(basePath+"assets", http.FS(&wrapAssetsFS{FS: assetsFS})) //3x-ui
	}
}
```
