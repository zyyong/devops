package main

import (
	"flag"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/api/action/artifactory"
	"github.com/yametech/devops/pkg/controller"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/devops", "127.0.0.1:3306")
	flag.Parse()

	store, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}

	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)
	//new artifactoryserver
	artifactory.NewArtifactoryServer("artifactory", server)
	//run artifactoryserver
	go func() {
		if err := server.Run("127.0.0.1:8080"); err != nil {
			errC <- err
		}
	}()

	go func() {
		if err := controller.NewWatchFlowRun(baseService).Run(); err != nil {
			errC <- err
		}
	}()

	if e := <-errC; e != nil {
		panic(e)
	}
}