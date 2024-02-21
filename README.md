# 3X-UI front-end separation To Windows


## Reproduce from 3x ui
## Temporarily migrate and store front-end files in the etc/html directory，The file does not contain front-end files.
## ! modify 

## modify

<details>
  <summary>Click for bug fix details</summary>

```
web/web.go
！del embed or //
////go:embed assets/*
//var assetsFS embed.FS

////go:embed html/*
//var htmlFS embed.FS

Compiling binary programs does not require including files in the HTML folder and assets folder
----------------------------------------------------------
! modify
type Server struct {
//
webDir string // dir
}
---------------------------------------------------------
func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	webDir := "/etc/www" // dir
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		webDir: webDir, //sing
	}
}
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


</details>


## goland : go build xxx.go .......finish complete
### Create a main.rc file with the following content: IDI_ ICON1 ICON "favicon. ico"
### Then run the terminal command in the current directory: windres - o main.syso main.rc
### Finally, go build completes packaging

## Recommended OS

- windows

## Languages

- English
- Farsi
- Chinese
- Russian
- Vietnamese
- Spanish

## bug fix

<details>
  <summary>Click for bug fix details</summary>

 
```
！Modify xray/process.go -->stop function
func (p *process) Stop() error {
	if !p.IsRunning() {
		return errors.New("xray is not running")
	}

	// 尝试发送 SIGTERM 信号
	err := p.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		// 如果发送 SIGTERM 失败，尝试直接强制终止进程
		err = p.cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("failed to stop xray: %v", err)
		}
	}

	// 等待进程退出
	_, err = p.cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for xray to exit: %v", err)
	}

	return nil
}

```
</details>


## bug fix for update

<details>
	
  <summary>Click for bug fix details</summary>

 
```
 -------new xray file-----config->config.go-------------

func GetXrayFolderPath() string {
	XrayFolderPath := os.Getenv("XUI_BIN_FOLDER")
	if XrayFolderPath == "" {
		XrayFolderPath = "/etc/xray"
	}
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(XrayFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(XrayFolderPath, os.ModePerm)
		if err != nil {
			// 处理创建目录失败的错误
			panic(err)
		}
	}
	return XrayFolderPath
}
func GetLogFolder() string {
	logFolderPath := os.Getenv("XUI_LOG_FOLDER")
	if logFolderPath == "" {
		logFolderPath = "/etc/log"
	}
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(logFolderPath, os.ModePerm)
		if err != nil {
			// 处理创建目录失败的错误
			panic(err)
		}
	}
	return logFolderPath
}
-------or web->web.go-----------------
	copyLinuxFiles := map[string]string{
		"xray":        xray.GetBinaryPath(),
		"geosite.dat": xray.GetGeositePath(),
		"geoip.dat":   xray.GetGeoipPath(),
	}
	copyWinFiles := map[string]string{
		"xray.exe":    xray.GetBinaryPath(),
		"wxray.exe":   xray.GetWxraytPath(),
		"geosite.dat": xray.GetGeositePath(),
		"geoip.dat":   xray.GetGeoipPath(),
	}
	if runtime.GOOS == "linux" {
		for fileName, filePath := range copyLinuxFiles {
			err := copyZipFile(fileName, filePath)
			if err != nil {
				return err
			}
		}
	} else if runtime.GOOS == "windows" {
		for fileName, filePath := range copyWinFiles {
			err := copyZipFile(fileName, filePath)
			if err != nil {
				return err
			}
		}
	} else {
		//return fmt.Errorf("不支持的操作系统：%s", runtime.GOOS)
		fmt.Errorf("不支持的操作系统：s%", runtime.GOOS)
	}
-----------------Check the xray operation every 5 seconds---------------------------
func (s *Server) startTask() {
	err := s.xrayService.RestartXray(true)
	if err != nil {
		logger.Warning("start xray failed:", err)
	}
	// Check whether xray is running every second
	s.cron.AddJob("@every 5s", job.NewCheckXrayRunningJob())
-----------------------xray->process.go for windows --------------------------------

// -------new way for windows or linux-----
func GetBinaryName() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("xray-%s-%s.exe", runtime.GOOS, runtime.GOARCH)
	}
	return fmt.Sprintf("xray-%s-%s", runtime.GOOS, runtime.GOARCH)
}

func GetBinaryPath() string {
	return config.GetXrayFolderPath() + "/" + GetBinaryName()
}

func GetConfigPath() string {
	return config.GetBinFolderPath() + "/config.json"
}

// -----file move to /etc/xray
func GetWxraytPath() string {
	return config.GetXrayFolderPath() + "/" + "wxray.exe"
}
--------------------------------------------------------------------------------------------------

！Modify web/service/server.go -->update function
	// 根据操作系统选择性地复制文件
	switch runtime.GOOS {
	case "windows":
		if err := copyZipFile("xray.exe", xray.GetBinaryPath()); err != nil {
			return err
		}
		if err := copyZipFile("wxray.exe", xray.GetWxraytPath()); err != nil {
			return err
		}
		if err := copyZipFile("geosite.dat", xray.GetGeositePath()); err != nil {
			return err
		}
		if err := copyZipFile("geoip.dat", xray.GetGeoipPath()); err != nil {
			return err
		}
	case "linux":
		if err := copyZipFile("xray", xray.GetBinaryPath()); err != nil {
			return err
		}
		if err := copyZipFile("geosite.dat", xray.GetGeositePath()); err != nil {
			return err
		}
		if err := copyZipFile("geoip.dat", xray.GetGeoipPath()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支持的操作系统：%s", runtime.GOOS)
	}

```

</details>

## Default Settings

<details>
  <summary>Click for default settings details</summary>

  ### Information

- **Port:** 2053
- **Username & Password:** It will be generated randomly if you skip modifying.
- **Database Path:**
  - /etc/x-ui/x-ui.db
- **Xray Config Path:**
  - /etc/config.json
- **Web Panel Path w/o Deploying SSL:**
  - http://ip:2053/panel
  - http://domain:2053/panel
- **Web Panel Path w/ Deploying SSL:**
  - https://domain:2053/panel
 
</details>

## API Routes

<details>
  <summary>Click for API routes details</summary>

#### Usage

- `/login` with `POST` user data: `{username: '', password: ''}` for login
- `/panel/api/inbounds` base for following actions:

| Method | Path                               | Action                                      |
| :----: | ---------------------------------- | ------------------------------------------- |
| `GET`  | `"/list"`                          | Get all inbounds                            |
| `GET`  | `"/get/:id"`                       | Get inbound with inbound.id                 |
| `GET`  | `"/getClientTraffics/:email"`      | Get Client Traffics with email              |
| `GET`  | `"/createbackup"`                  | Telegram bot sends backup to admins         |
| `POST` | `"/add"`                           | Add inbound                                 |
| `POST` | `"/del/:id"`                       | Delete Inbound                              |
| `POST` | `"/update/:id"`                    | Update Inbound                              |
| `POST` | `"/clientIps/:email"`              | Client Ip address                           |
| `POST` | `"/clearClientIps/:email"`         | Clear Client Ip address                     |
| `POST` | `"/addClient"`                     | Add Client to inbound                       |
| `POST` | `"/:id/delClient/:clientId"`       | Delete Client by clientId\*                 |
| `POST` | `"/updateClient/:clientId"`        | Update Client by clientId\*                 |
| `POST` | `"/:id/resetClientTraffic/:email"` | Reset Client's Traffic                      |
| `POST` | `"/resetAllTraffics"`              | Reset traffics of all inbounds              |
| `POST` | `"/resetAllClientTraffics/:id"`    | Reset traffics of all clients in an inbound |
| `POST` | `"/delDepletedClients/:id"`        | Delete inbound depleted clients (-1: all)   |
| `POST` | `"/onlines"`                       | Get Online users ( list of emails )       |

\*- The field `clientId` should be filled by:

- `client.id` for VMESS and VLESS
- `client.password` for TROJAN
- `client.email` for Shadowsocks


- [API Documentation](https://documenter.getpostman.com/view/16802678/2s9YkgD5jm)
- [<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/16802678-1a4c9270-ac77-40ed-959a-7aa56dc4a415?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D16802678-1a4c9270-ac77-40ed-959a-7aa56dc4a415%26entityType%3Dcollection%26workspaceId%3D2cd38c01-c851-4a15-a972-f181c23359d9)
</details>

## Environment Variables

<details>
  <summary>Click for environment variables details</summary>

#### Usage

| Variable       |                      Type                      | Default       |
| -------------- | :--------------------------------------------: | :------------ |
| XUI_LOG_LEVEL  | `"debug"` \| `"info"` \| `"warn"` \| `"error"` | `"info"`      |
| XUI_DEBUG      |                   `boolean`                    | `false`       |
| XUI_BIN_FOLDER |                    `string`                    | `"bin"`       |
| XUI_DB_FOLDER  |                    `string`                    | `"/etc/x-ui"` |
| XUI_LOG_FOLDER |                    `string`                    | `"/etc/log"`  |

Example:

```sh
XUI_BIN_FOLDER="bin" XUI_DB_FOLDER="/etc/x-ui" go build main.go
```

</details>
