package cmd

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/log"
	//"github.com/will835559313/apiman/pkg/myvalidator"
	"github.com/will835559313/apiman/pkg/jwt"
	"github.com/will835559313/apiman/pkg/setting"
	"github.com/will835559313/apiman/routes"
	"github.com/will835559313/apiman/routes/user"
	"gopkg.in/go-playground/validator.v9"
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func newWeb() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
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

	// get db connection
	models.Dbinit()

	// migrate tables
	models.DbMigrate()

	// set jwt
	jwt.JwtInint()

	address := ":" + port
	e := newWeb()
	e.GET("/", routes.Index)
	//e.GET("/home", routes.Home)
	e.POST("/users", user.CreateUser)
	e.GET("/users/:username", user.GetUserByName)
	e.PUT("/users/:username", user.UpdateUserByName)
	e.DELETE("/users/:username", user.DeleteUserByName)
	e.POST("/users/:username/reset_password", user.RestUserPassword)

	e.POST("/oauth2/token", user.GetToken)

	e.Logger.Fatal(e.Start(address))
	return nil
}
