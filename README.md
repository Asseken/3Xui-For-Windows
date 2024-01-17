# 3X-UI front-end separation

## Temporarily migrate and store front-end files in the/etc/wwwww directory
## ! modify

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
## goland : go build xxx.go .......finish complete

## Recommended OS

- Ubuntu 20.04+
- Debian 11+
- CentOS 8+
- Fedora 36+
- Arch Linux
- Manjaro
- Armbian
- AlmaLinux 9+
- Rockylinux 9+


## Languages

- English
- Farsi
- Chinese
- Russian
- Vietnamese
- Spanish

## Default Settings

<details>
  <summary>Click for default settings details</summary>

  ### Information

- **Port:** 2053
- **Username & Password:** It will be generated randomly if you skip modifying.
- **Database Path:**
  - /etc/x-ui/x-ui.db
- **Xray Config Path:**
  - /usr/local/x-ui/bin/config.json
- **Web Panel Path w/o Deploying SSL:**
  - http://ip:2053/panel
  - http://domain:2053/panel
- **Web Panel Path w/ Deploying SSL:**
  - https://domain:2053/panel
 
</details>

## IP Limit

<details>
  <summary>Click for IP limit details</summary>

#### Usage

**Note:** IP Limit won't work correctly when using IP Tunnel

- For versions up to `v1.6.1`:

  - IP limit is built-in into the panel.

- For versions `v1.7.0` and newer:

  - To make IP Limit work properly, you need to install fail2ban and its required files by following these steps:

    1. Use the `x-ui` command inside the shell.
    2. Select `IP Limit Management`.
    3. Choose the appropriate options based on your needs.
   
  - make sure you have access.log on your Xray Configuration
  
  ```sh
    "log": {
    "loglevel": "warning",
    "access": "./access.log",
    "error": "./error.log"
    },
  ```

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
