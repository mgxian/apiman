package cmd

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	mylog "github.com/will835559313/apiman/pkg/log"
	"github.com/will835559313/apiman/pkg/setting"
	"github.com/will835559313/apiman/routes"
	"github.com/will835559313/apiman/routes/api"
	"github.com/will835559313/apiman/routes/apigroup"
	"github.com/will835559313/apiman/routes/project"
	"github.com/will835559313/apiman/routes/team"
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
	mylog.LoggerInit()

	// get db connection
	models.Dbinit()

	// migrate tables
	models.DbMigrate()

	// set jwt
	jwt.JwtInint()

	address := ":" + port
	e := newWeb()
	e.GET("/", routes.Index)

	// user
	e.POST("/users", user.CreateUser)
	e.GET("/users/:username", user.GetUserByName)
	e.PUT("/users/:username", user.UpdateUserByName)
	e.DELETE("/users/:username", user.DeleteUserByName)
	e.POST("/users/:username/change_password", user.ChangeUserPassword)
	e.GET("/users/:username/teams", user.GetUserTeams)
	e.GET("/users/:username/projects", user.GetUserProjects)

	// token
	e.POST("/oauth2/token", user.GetToken)

	// team
	e.POST("/teams", team.CreateTeam)
	e.GET("/teams/:teamname", team.GetTeamByName)
	e.PUT("/teams/:teamname", team.UpdateTeamByName)
	e.DELETE("/teams/:teamname", team.DeleteTeamByName)
	e.POST("/teams/:teamname/members", team.AddOrUpdateTeamMember)
	e.DELETE("/teams/:teamname/members/:username", team.RemoveTeamMember)
	e.PUT("/teams/:teamname/members/:username", team.AddOrUpdateTeamMember)
	e.GET("/teams/:teamname/members", team.GetTeamMembers)
	e.GET("/teams/:teamname/members/:username", team.GetTeamMember)
	e.GET("/teams/:teamname/projects", team.GetTeamProjets)

	// project
	e.POST("/teams/:teamname/projects", project.CreateProject)
	e.GET("/projects/:id", project.GetProjectByID)
	e.PUT("/projects/:id", project.UpdateProjectByID)
	e.DELETE("/projects/:id", project.DeleteProjectByID)
	e.POST("/projects/:id/migrate", project.MigrateProjectByID)
	e.GET("/projects/:id/apis", project.GetProjectApis)
	e.GET("/projects/:id/apigroups", project.GetProjectApiGroups)

	// apigroup
	e.POST("/projects/:id/apigroups", apigroup.CreateApiGroup)
	e.GET("/apigroups/:id", apigroup.GetApiGroupByID)
	e.PUT("/apigroups/:id", apigroup.UpdateApiGroupByID)
	e.DELETE("/apigroups/:id", apigroup.DeleteApiGroupByID)
	e.GET("/apigroups/:id/apis", apigroup.GetApiGroupApis)

	// api
	e.POST("/apigroups/:id/apis", api.CreateApi)
	e.POST("/projects/:id/apis", api.CreateDefaultApi)
	e.GET("/apis/:id", api.GetApi)
	e.PUT("/apis/:id", api.UpdateApi)
	e.DELETE("/apis/:id", api.DeleteApi)

	// log start
	e.Logger.Fatal(e.Start(address))

	return nil
}
