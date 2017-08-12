package main

import (
	"fmt"
	"net/http"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
)

type args struct {
	Port int    `cli:"p,port" usage:"port to bind to"`
	Host string `cli:"h,host" usage:"host to bind to"`
	Dir  string `cli:"d,directory" usage:"directory to serve from"`
}

var (
	host string
	port int
	dir  string
)

func main() {
	cli.Run(new(args), func(ctx *cli.Context) error {
		a := ctx.Argv().(*args)
		host = a.Host
		if host == "" {
			host = "0.0.0.0"
		}
		port = a.Port
		if port == 0 {
			port = 8080
		}
		dir = a.Dir
		if dir == "" {
			a, err := filepath.Abs(".")
			if err != nil {
				fatalf("Could not find current directory: %v", err)
			}
			dir = a
		}
		return nil
	})
	printf("Serving files on %s:%d from %q\n", host, port, mustAbs(dir))
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.StaticFS("./", http.FileSystem(http.Dir(dir)))
	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		fatalf("Error occured while serving files: %v", err)
	}
}

func mustAbs(relpath string) string {
	abs, err := filepath.Abs(relpath)
	if err != nil {
		fatalf("Could not resolve absolute path: %v", err)
	}
	return abs
}

func printf(tmpl string, values ...interface{}) {
	fmt.Printf(tmpl, values...)
}

func fatalf(tmpl string, values ...interface{}) {
	fmt.Printf(tmpl+"\n", values...)
	os.Exit(2)
}
