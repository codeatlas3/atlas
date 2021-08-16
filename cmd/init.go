/*
Copyright © 2021 Mohanasundaram N <codewithmohanasundaram@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var Dir string
var Name string
var Lang string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize new project",
	Long:    `Initialize new Golang/Flutter project and the current directory`,
	Example: "atalas init",
	Run: func(cmd *cobra.Command, args []string) {
		projectSelector := promptui.Select{
			Label: "Select the language",
			Items: []string{"Golang", "Dart(Flutter)"},
		}

		_, selected, err := projectSelector.Run()

		if err != nil {
			log.Fatalln("Unable to select project" + err.Error())
		}

		projectNamePrompt := promptui.Prompt{
			Label: "What is the name of the project?",
			Validate: func(s string) error {
				if s != "" && len(s) < 2 {
					return errors.New("Project must be greater than 2 characters ")
				}
				return nil
			},
		}

		projectName, err := projectNamePrompt.Run()

		if err != nil {
			log.Fatalln("Prompt failed" + err.Error())
		}

		println(fmt.Sprintf("🗂️ Initializing %s project", selected))

		switch selected {
		case "Golang":
			createGolangProject(projectName)
			break
		case "Dart(Flutter)":
			createFlutterProject(projectName)
			break
		default:
			break
		}

		println("✔ " + projectName + " initialized successfully")

		runIDESelect := promptui.Select{
			Label: "Do you want want open the generated project in Golang?",
			Items: []string{"Y", "N"},
		}

		_, val, err := runIDESelect.Run()

		if val == "Y" {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalln("Unable to find the working directory")
			}
			cmd := exec.Command("goland", wd)
			err = cmd.Run()

			if err != nil {
				log.Fatalln("Unable to launch project in goland")
			}
		}

	},
}

func createFlutterProject(projectName string) {
	println(projectName)
}

func createGolangProject(projectName string) {
	err := os.MkdirAll(projectName, os.ModePerm)
	if err != nil {
		log.Fatalln("Unable to create directory")
	}
	err = os.Chdir(projectName)

	if err != nil {
		log.Fatalln("Unable to change directory")
	}

	createGOMODFile(projectName)
	createMainGoFile(projectName)
	createServerGoFile(projectName)
	createDatabaseGoFile(projectName)
	createRouteGoFile(projectName)
	createHandlersGoFile(projectName)
	createRepositoryGoFile(projectName)
	createModelGoFile(projectName)
	createAPIGoFile(projectName)
	createLoggyGoFile(projectName)
	createEnvGoFile(projectName)
}

func createEnvGoFile(projectName string) {
	str := `package utils

	const LOCAL = "local"
	const DEV = "dev"
	const TESTING = "testing"
	const PROD = "prod"
`
	createFileFromTemplate("./utils/environment.go", str, projectName)
}

func createLoggyGoFile(projectName string) {
	str := `package loggy

	import (
		"go.uber.org/zap"
	)
	
	type Loggy interface {
		Info(interface{})
		Error(interface{})
		Warn(interface{})
		Debug(interface{})
	}
	
	func NewLoggy(name string) Loggy {
		logger, err := zap.NewProduction()
		if err != nil {
			println("unable to init logger")
		}
		defer logger.Sync()
		logger.Named(name)
		return loggy{
			logger: logger,
		}
	}
	
	type loggy struct {
		logger *zap.Logger
	}
	
	func (l loggy) Info(v interface{}) {
		l.logger.Sugar().Info(v)
	}
	
	func (l loggy) Debug(v interface{}) {
		l.logger.Sugar().Debug(v)
	}
	
	func (l loggy) Warn(v interface{}) {
		l.logger.Sugar().Warn(v)
	}
	
	func (l loggy) Error(v interface{}) {
		l.logger.Sugar().Error(v)
	}
`
	err := os.MkdirAll("utils/loggy", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create utils/loggy directory")
	}
	createFileFromTemplate("./utils/loggy/loggy.go", str, projectName)
}

func createAPIGoFile(projectName string) {
	str := `package api

const StatusInternalServerError = "Something went wrong, please try after sometimes"
const StatusBadRequest = "Bad request"
`
	err := os.MkdirAll("utils/api", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create utils/api directory")
	}
	createFileFromTemplate("./utils/api/error_message.go", str, projectName)
}

func createModelGoFile(projectName string) {
	str := `package models

type SuccessResponse struct {
	Success bool       	json:"success"
	Data    interface{} json:"data"
}

type FailureResponse struct {
	Success bool        json:"success"
	Message string      json:"message"
	Error   interface{} json:"error"
}
`
	err := os.MkdirAll("model", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create model directory")
	}
	createFileFromTemplate("./model/response.go", str, projectName)
}

func createRepositoryGoFile(projectName string) {
	str := `package repository
	
	import (
		"context"
		models "{{ .ProjectName }}/model"
		"{{ .ProjectName }}/utils/loggy"
		"time"
	
		"go.mongodb.org/mongo-driver/bson"
	
		"go.mongodb.org/mongo-driver/mongo"
	
		"go.mongodb.org/mongo-driver/bson/primitive"
	)
	
	type Repository interface {}

	func NewRepository(collections database.Collections, loggy loggy.Loggy) Repository {
	return repository{
		collections: collections,
		loggy:       loggy,
	}
}
`

	err := os.MkdirAll("repository", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create repository directory")
	}
	createFileFromTemplate("./repository/repository.go", str, projectName)
}

func createHandlersGoFile(projectName string) {
	str := `package handler

	import (
		models "{{ .ProjectName }}/model"
		"{{ .ProjectName }}/repository"
		"{{ .ProjectName }}/utils/api"
		"{{ .ProjectName }}/utils/loggy"
	
		"go.mongodb.org/mongo-driver/bson/primitive"
	
		"github.com/gofiber/fiber/v2"
	)
	
	// Handler
	type Handler interface {}
	
	func NewHandler(repository repository.Repository, loggy loggy.Loggy) Handler {
		return handler{
			repository: repository,
			loggy:      loggy,
		}
	}
`
	err := os.MkdirAll("handler", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create handler directory")
	}
	createFileFromTemplate("./handler/handler.go", str, projectName)
}

func createRouteGoFile(projectName string) {
	str := `package router

	import (
		"{{ .ProjectName }}/handler"
	
		"github.com/gofiber/fiber/v2"
	)
	
	func SetupRoutes(app *fiber.App, h handler.Handler) {
	
		app.Get("/health-check", func(c *fiber.Ctx) error {
			return c.Status(200).JSON(fiber.Map{"success": true, "message": "welcome"})
		})
	
		v1 := app.Group("/api/v1")
	}
`
	err := os.MkdirAll("router", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create router directory")
	}
	createFileFromTemplate("./router/router.go", str, projectName)
}

func createDatabaseGoFile(projectName string) {
	str := `package database

	import (
		"context"
		"{{ .ProjectName }}/utils/loggy"
		"os"
		"time"
	
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"
		"go.mongodb.org/mongo-driver/mongo/readpref"
	)

	type Collections struct {}
	
	func Connect(loggy loggy.Loggy) (Collections, error) {
		uri := os.Getenv("MONGO_URI")
		dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		opt := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(dbCtx, opt)
		loggy.Error(err)
		if err != nil {
			return Collections{}, err
		} else {
			if err := client.Ping(dbCtx, readpref.Primary()); err != nil {
				loggy.Error(err)
				return Collections{}, err
			}
			// Replace with your db name
			db := client.Database("test")
			
			sampleCollection := db.Collection("sample")

			return Collections{}, nil
		}
	}
`
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		log.Fatalln("unable to create database directory")
	}
	createFileFromTemplate("./database/database.go", str, projectName)
}

func createServerGoFile(projectName string) {
	str := `package server

	import (
		"{{ .ProjectName }}/database"
		"{{ .ProjectName }}/handler"
		"{{ .ProjectName }}/repository"
		"{{ .ProjectName }}/router"
		"{{ .ProjectName }}/utils/loggy"
		"log"
	
		"github.com/gofiber/fiber/v2"
		"github.com/gofiber/fiber/v2/middleware/recover"
	)
	
	func Setup() *fiber.App {
		l := loggy.NewLoggy("{{ .ProjectName }}")
	
		collections, err := database.Connect(l)
	
		if err != nil {
			log.Fatalln("Unable connect to database: " + err.Error())
		}
	
		app := fiber.New()
	
		app.Use(recover.New())
	
		r := repository.NewRepository(collections, l)
	
		h := handler.NewHandler(r, l)
	
		router.SetupRoutes(app, h)
	
		return app
	}
	`
	err := os.MkdirAll("server", os.ModePerm)
	if err != nil {
		log.Fatalln("Unable to create server directory")
	}
	createFileFromTemplate("./server/server.go", str, projectName)
}

func createGOMODFile(projectName string) {
	str := `module {{ .ProjectName }}

	go 1.16
	
	require (
		github.com/andybalholm/brotli v1.0.3 // indirect
		github.com/fatih/structs v1.1.0 // indirect
		github.com/gavv/httpexpect/v2 v2.3.1 // indirect
		github.com/gofiber/fiber/v2 v2.15.0 // indirect
		github.com/golang/snappy v0.0.4 // indirect
		github.com/google/go-querystring v1.1.0 // indirect
		github.com/imkira/go-interpol v1.1.0 // indirect
		github.com/klauspost/compress v1.13.1 // indirect
		github.com/sergi/go-diff v1.2.0 // indirect
		github.com/stretchr/objx v0.3.0 // indirect
		github.com/stretchr/testify v1.7.0 // indirect
		github.com/valyala/fasthttp v1.28.0 // indirect
		github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
		github.com/xeipuuv/gojsonschema v1.2.0 // indirect
		github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
		go.mongodb.org/mongo-driver v1.6.0 // indirect
		go.uber.org/atomic v1.9.0 // indirect
		go.uber.org/multierr v1.7.0 // indirect
		go.uber.org/zap v1.18.1 // indirect
		golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
		golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
		golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
		golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
		gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	)
`
	createFileFromTemplate("./go.mod", str, projectName)
}

func createMainGoFile(projectName string) {
	str := `package main

	import (
		"log"
		"os"
		"vehicle_management_service/server"
		"vehicle_management_service/utils"
	)
	
	func main() {
		port := os.Getenv("PORT")
		env := os.Getenv("ENV")
		app := server.Setup()
		if env == utils.LOCAL {
			log.Fatal(app.Listen("localhost:" + port))
		} else {
			log.Fatal(app.Listen(":" + port))
		}
	}

	`
	createFileFromTemplate("./main.go", str, projectName)

}

func createFileFromTemplate(path string, str string, projectName string) {
	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	t := template.Must(template.New("").Parse(str))

	err = t.Execute(f, map[string]interface{}{
		"ProjectName": projectName,
	})

	if err != nil {
		log.Fatalf("\n %v template generation filed", path)
	}

	err = f.Close()

	if err != nil {
		log.Fatalln("unable to close file")
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	initCmd.PersistentFlags().StringVarP(&Dir, "directory", "d", "", "Directory where the generated project lives")
	initCmd.PersistentFlags().StringVarP(&Name, "name", "n", "", "Name of the project to be generated")
	initCmd.PersistentFlags().StringVarP(&Lang, "language", "l", "", "Either Golang or Flutter) ")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
