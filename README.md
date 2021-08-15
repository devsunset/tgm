# tgm

	########## dev ##########
	go install & path setting -> https://golang.org/
	EchoHigh performance, extensible, minimalist Go web framework
	https://echo.labstack.com/

	mkdir -p <work_path>/tgm
	cd <work_path>/tgm

	go mod init tgm
	go get github.com/labstack/echo/v4
	go get github.com/labstack/echo/middleware
	go mod download golang.org/x/time
    go get github.com/swaggo/swag/cmd/swag
    go get -u github.com/swaggo/echo-swagger
	go mod download github.com/alecthomas/template
	go get github.com/alecthomas/template@v0.0.0-20190718012654-fb15b899a751

	https://pkg.go.dev/github.com/amarnus/echo-swagger#readme-echo-swagger
    https://github.com/swaggo/swag#declarative-comments-format
