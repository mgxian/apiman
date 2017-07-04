// apiman project apiman.go
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/will835559313/apiman/cmd"
)

const APP_VER = "0.0.1"

func main() {
	fmt.Println(APP_VER)
	app := cli.NewApp()
	app.Name = "apiman"
	app.Usage = "api manager"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
