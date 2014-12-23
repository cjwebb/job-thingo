package main

import (
	"fmt"
	"log"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html"
	"html/template"
	"github.com/cjwebb/job-thingo/db"
)

type JobLink struct {
	id    string
	email string
}

func main() {
	database := db.NewDatabase()
	startApp(database)
}

func startApp(db db.Database) {

	helpers := template.FuncMap{
		"unescape": func(s string) template.HTML {
			return template.HTML(html.UnescapeString(s))
		},
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "base",
		Funcs:  []template.FuncMap{helpers},
	}))

	// Home page
	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", nil)
	})

	// Display form to create new JobLink
	m.Get("/gen", func(params martini.Params, r render.Render) {
		r.HTML(200, "generate", nil)
	})

	// Generate new JobLink and redirect
	m.Post("/gen", func(params martini.Params, r render.Render) {
		fmt.Print(params)
		// https://github.com/martini-contrib/binding
		r.HTML(200, "job", nil)
	})

	// Show incoming JobLink
	m.Get("/a/:jobid", func(params martini.Params, r render.Render) {
		response, err := db.GetJob(params["jobid"])

		if err != nil {
			log.Print(err.Error())
			r.Redirect("/")
		} else {
			r.HTML(200, "job", response)
		}
	})

	// About page
	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", nil)
	})

	m.Run()
}
