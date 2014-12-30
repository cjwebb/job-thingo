package main

import (
	"fmt"
	"log"
	"html"
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"

	"github.com/satori/go.uuid"

	"github.com/cjwebb/job-thingo/db"
)

type JobForm struct {
	Title		string `form:"title"       binding:"required"`
	Description	string `form:"description" binding:"required"`
	Rate		string `form:"rate"        binding:"required"`
	Email		string `form:"email"       binding:"required"`
}

func main() {
	database := db.NewDatabase()
	startApp(database)
}

func startApp(database db.Database) {

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

	// Generate new Job and redirect
	m.Post("/gen-job", binding.Form(JobForm{}), func(form JobForm, r render.Render, errs binding.Errors) {
		if errs != nil {
			for _, e := range errs {
				fmt.Print(e.FieldNames[0])
				fmt.Println(e.Message)
			}
			// todo(cjwebb) - santize form data
			r.HTML(400, "home", form)
		} else {
			u1 := uuid.NewV4().String()
			// todo(cjwebb) - sanitize form data
			// todo(cjwebb) - how should we handle errors here?
			database.PutJob(db.Job{u1, form.Title, form.Description, form.Email, form.Rate})
			r.Redirect("/a/"+u1)
		}
	})

	// Show incoming Job
	m.Get("/a/:jobid", func(params martini.Params, r render.Render) {
		jobId, err := uuid.FromString(params["jobid"])
		if err != nil {
			// input not a UUID, so don't try database lookup
			r.Redirect("/")
		} else {
			job, err := database.GetJob(jobId.String())
			if err != nil {
				log.Print(err.Error())
				r.Redirect("/")
			} else {
				r.HTML(200, "job", job)
			}
		}
	})

	// About page
	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", nil)
	})

	m.Run()
}
