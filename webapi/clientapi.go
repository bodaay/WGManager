package webapi

import (
	"WGManager/webapi/resource"
	"WGManager/wg"
	"bytes"
	"mime"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//StartAdminClient start the REST API Echo Server for inserting watermark
func StartClient(wgConfig *wg.WGConfig) error {
	e := echo.New()
	const subserviceIdentifier = "StartWebClient"
	configureClientWebServer(e)
	configureAllRoutesClient(e, wgConfig)
	address := (wgConfig.APIListenAddress + ":" + strconv.Itoa(int(wgConfig.APIListenPort)))
	//err := e.StartTLS(address, (config.RootCertFile), (config.RootCertKey))
	err := e.Start(address)
	if err != nil {
		panic("Error StartWebClient StartTLS")
	}
	return nil
}
func configureClientWebServer(e *echo.Echo) {
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

func configureAllRoutesClient(e *echo.Echo, wgConfig *wg.WGConfig) {
	postAllocateClient(e, wgConfig)
	postRevokeClient(e, wgConfig)
}

func postAllocateClient(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.POST("/api/client", func(c echo.Context) error {
		u := new(resource.WgAllocateClientRequest)
		if err := c.Bind(u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		qrbytes, err := wgConfig.AllocateClient(u.Instancename, u.Clientuuid)
		responseObj := "Allocation Successfull"
		if err != nil {
			responseObj = err.Error()
			return c.JSONPretty(http.StatusBadRequest, responseObj, "  ")
		}
		return c.Stream(http.StatusOK, "image/jpeg", bytes.NewReader(qrbytes))
		//return c.JSONPretty(http.StatusOK, responseObj, "  ")
	})
}

func postRevokeClient(e *echo.Echo, wgConfig *wg.WGConfig) {
	e.DELETE("/api/client", func(c echo.Context) error {
		u := new(resource.WgRevokeClientRequest)
		if err := c.Bind(u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		err := wgConfig.RevokeClient(u.Instancename, u.Clientuuid)
		responseObj := "Revocation Successfull"
		if err != nil {
			responseObj = err.Error()
			return c.JSONPretty(http.StatusBadRequest, responseObj, "  ")

		}

		return c.JSONPretty(http.StatusOK, responseObj, "  ")
	})
}
