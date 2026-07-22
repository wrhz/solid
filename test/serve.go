package test

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	solidConfig "github.com/wrhz/solid/config"
	solidInit "github.com/wrhz/solid/init"
	solidRoute "github.com/wrhz/solid/route"
)

var serve *mux.Router

func TestRoute(callFunc func()) error {
	workDir := os.Getenv("TEST_WORKDIR")
    if workDir == "" {
        return fmt.Errorf("TEST_WORKDIR not set")
    }

	os.Chdir(workDir)

	serverConfig := solidConfig.GetServerConfig()

	var mainStruct solidRoute.SolidMainRoute
	var err error

	mainStruct, serve, err = solidInit.InitServer(serverConfig)

	if err != nil {
		return err
	}

	if serve == nil {
       return fmt.Errorf("serve is nil, did you forget to call SetServe in TestMain?")
    }

	httpServer := serverConfig.GetServerConfig()

	httpServer.Addr = ":" + strconv.Itoa(serverConfig.GetPort())
	httpServer.Handler = serve
	httpServer.TLSConfig = serverConfig.GetTLSConfig()

	callFunc()

	if err := mainStruct.ServerEnd(); err != nil {
		return err
	}

	return nil
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.ServeHTTP(w, r)
}