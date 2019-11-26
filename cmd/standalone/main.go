package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lukemgriffith/captainhook"
	"github.com/lukemgriffith/captainhook/configparser"
	"github.com/lukemgriffith/captainhook/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


func main() {

	var configPath, secretPath, passphrase string

	flag.StringVar(&configPath, "configPath", "config.yml", "YAML file to configure the service with endpoints.")
	flag.StringVar (&secretPath, "secretPath", "", "Encrypted YAML blob containing string map of secrets.")
	flag.StringVar (&passphrase, "passphrase", "", "Passphrase for encrypted YAML blob.")
	flag.Parse()

	var secEnd captainhook.SecretEngine = createSecretsEngine(secretPath, passphrase)

	svc, err := configparser.NewConfig(configPath)

	if err != nil {
		panic(err)
	}


	endpoints, err := svc.Endpoints()

	if err != nil {
		panic(err)
	}

	if secEnd != nil {
		// Validating all configured secrets are available.
		for _, end := range endpoints {
			for _, name := range end.Secrets {
				_, err := secEnd.GetTextSecret(name)

				if err != nil {
					panic(fmt.Sprintf("unable to find secret %s", name))
				}
			}
		}
	}


	app := server.New(svc, secEnd)

	hookserver := &http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	go hookserver.ListenAndServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	_ = <-c
	log.Print("os signal recieved processing.")

	log.Print("shutting server down gracefully.")
	hookserver.Shutdown(context.Background())

}

func createSecretsEngine(secretPath string, passphrase string) captainhook.SecretEngine {
	if secretPath != "" && passphrase != "" {
		secEnd, err := configparser.NewSecretEngine(secretPath, passphrase)

		if err != nil {
			panic(err)
		}

		return secEnd
	} else {
		return nil
	}
}
