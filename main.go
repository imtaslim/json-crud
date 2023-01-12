package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/schema"
	"github.com/spf13/afero"
	"github.com/yookoala/realpath"
	
	"github.com/imtaslim/json-crud/handler"
)

func main (){
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	
	assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		log.Panic(err)
	}
	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))
	
	h := handler.New(asst, decoder)
	if err := http.ListenAndServe("localhost:3000", h); err != nil {
		log.Panic(err)
	}
}