package webapi

import (
	"WGManager/webapi/resource"
	"WGManager/wg"
	"mime"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//StartAdminClient start the REST API Echo Server for inserting watermark
func StartAdminClient(wgConfig *wg.WGConfig) error {
	e := echo.New()
	const subserviceIdentifier = "StartWebClient"
	configureAdminClientWebServer(e)
	configureAdminAllRoutesClient(e, wgConfig)
	address := (wgConfig.AdminAPIListenAddress + ":" + strconv.Itoa(int(wgConfig.AdminAPIListenPort)))
	//err := e.StartTLS(address, (config.RootCertFile), (config.RootCertKey))
	err := e.Start(address)
	if err != nil {
		panic("Error StartWebClient StartTLS")
	}
	return nil
}
func configureAdminClientWebServer(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("100M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	mime.AddExtensionType(".js", "application/javascript") //This will solve some windows shit issue, when it will serve javascript file as text/plain, read more about it at:https://github.com/labstack/echo/issues/1038

}

func configureAdminAllRoutesClient(e *echo.Echo, wgConfig *wg.WGConfig) {
	getAllWgInstance(e, wgConfig)
	postAddWgInstance(e, wgConfig)
	postRemoveWgInstance(e, wgConfig)
	postDeployWgInstance(e, wgConfig)
}
func getAllWgInstance(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.GET("/api/instance", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, wgConfig.WGInstances, "  ")
	})
}
func postAddWgInstance(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.PUT("/api/instance", func(c echo.Context) error {
		u := new(resource.WgAddInstanceRequest)
		if err := c.Bind(u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		err := wgConfig.CreateNewInstance(u.IP, u.Port, u.DNS, u.UseNAT, u.EthernetName, uint64(u.MaxClient))
		responseObj := "Add Successfull"
		if err != nil {
			responseObj = err.Error()
			return c.JSONPretty(http.StatusBadRequest, responseObj, "  ")

		}
		return c.JSONPretty(http.StatusOK, responseObj, "  ")
	})
}
func postRemoveWgInstance(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.DELETE("/api/instance", func(c echo.Context) error {
		u := new(resource.WgRemoveInstanceRequest)
		if err := c.Bind(u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		err := wgConfig.RemoveInstanceByName(u.Instancename)
		responseObj := "Remove Successfull"
		if err != nil {
			responseObj = err.Error()
			return c.JSONPretty(http.StatusBadRequest, responseObj, "  ")

		}
		return c.JSONPretty(http.StatusOK, responseObj, "  ")
	})
}

func postDeployWgInstance(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.POST("/api/instance", func(c echo.Context) error {
		u := new(resource.WgDeploynstanceRequest)
		if err := c.Bind(u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		err := wgConfig.DeployInstanceByName(u.Instancename)
		responseObj := "Deploy Successfull"
		if err != nil {
			responseObj = err.Error()
			return c.JSONPretty(http.StatusBadRequest, responseObj, "  ")

		}

		return c.JSONPretty(http.StatusOK, responseObj, "  ")
	})
}
