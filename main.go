package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
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
	printf("Serving files from %q,\n", mustAbs(dir))
	addrs, err := ListAddress()
	if err != nil {
		fatalf("Error occured %v", err)
	}
	printf("Serving on:\n")
	for _, ip := range addrs {
		printf("-> %s\n", ip.String())
	}
	printf("\n")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.StaticFS("./", http.FileSystem(http.Dir(dir)))
	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		fatalf("Error occured while serving files: %v", err)
	}
}

// ListAddress returns a list of ip addresses found across all network
// interfaces.
func ListAddress() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "reading network interfaces")
	}
	ret := []net.IP{}
	for _, i := range ifaces {
		addresses, err := i.Addrs()
		if err != nil {
			return nil, errors.Wrapf(err, "reading ip address of %v", litter.Sdump(i))
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ret = append(ret, ip)
		}
	}
	return ret, nil
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
