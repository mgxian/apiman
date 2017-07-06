package cmd

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"github.com/will835559313/apiman/pkg/log"
	"github.com/will835559313/apiman/pkg/setting"
	"github.com/will835559313/apiman/routes"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: `Apiman web server`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "conf/app.conf", "config file path"),
	},
}

func newWeb() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func runWeb(c *cli.Context) error {
	port := "5000"
	if c.IsSet("port") {
		port = c.String("port")
	}

	if c.IsSet("config") {
		setting.CustomConf = c.String("config")
	}

	// load config
	setting.NewConfig()

	// set logger
	log.LoggerInit()

	address := ":" + port
	e := newWeb()
	e.GET("/", routes.Index)
	e.GET("/home", routes.Home)

	e.Logger.Fatal(e.Start(address))
	return nil
}
