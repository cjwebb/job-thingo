package main

import (
	"fmt"
	"html"
	"html/template"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"github.com/microcosm-cc/bluemonday"
	"github.com/satori/go.uuid"

	"github.com/cjwebb/job-thingo/db"
)

type JobForm struct {
	Title		string `form:"title"       	binding:"required"`
	Description	string `form:"description"	binding:"required"`
	Rate		string `form:"rate"        	binding:"required"`
	ContactEmail	string `form:"contact_email"	binding:"required"`
	UserEmail	string `form:"user_email"  	binding:"required"`
}

type EmailForm struct {
	Email	string `form:"email"	binding:"required"`
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

	sanitizer := bluemonday.UGCPolicy()
	sanitizeForm := func(form JobForm) JobForm {
		return JobForm{
			sanitizer.Sanitize(form.Title),
			sanitizer.Sanitize(form.Description),
			sanitizer.Sanitize(form.Rate),
			sanitizer.Sanitize(form.ContactEmail),
			sanitizer.Sanitize(form.UserEmail),
		}
	}
	sanitizeEmailForm := func(form EmailForm) EmailForm {
		return EmailForm{
			sanitizer.Sanitize(form.Email),
		}
	}
	newId := func() string {
		return uuid.NewV4().String()
	}

	// Home page
	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", nil)
	})

	// Generate new Job and redirect
	m.Post("/gen-job", binding.Form(JobForm{}), func(form JobForm, r render.Render, errs binding.Errors) {
		saneForm := sanitizeForm(form)
		if errs != nil {
			// todo(cjwebb) - display errors
			r.HTML(400, "home", saneForm)
		} else {
			id := newId()
			genJob := db.Job{
				id,
				saneForm.Title,
				saneForm.Description,
				saneForm.ContactEmail,
				saneForm.Rate,
				[]db.JobRef{db.JobRef{id,saneForm.UserEmail}},
			}
			// todo(cjwebb) - how should we handle errors here?
			database.PutJob(genJob)
			r.Redirect("/a/" + id)
		}
	})

	// Show Job
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

	// Generate new link to Job
	m.Post("/a/:jobid/gen-link", binding.Form(EmailForm{}), func(form EmailForm, params martini.Params, r render.Render, errs binding.Errors) {
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
				saneForm := sanitizeEmailForm(form)
				fmt.Println(saneForm.Email)
				if errs != nil {
					// todo(cjwebb) - show errors
					r.HTML(400, "job", job)
				} else {
					// job exists, form is valid... so gen new link
					id := newId()
					newRef := db.JobRef{id, saneForm.Email}
					// todo(cjwebb) - how should we handle errors here?
					genJob := db.Job{
						id,
						job.Title,
						job.Description,
						job.ContactEmail,
						job.Rate,
						append([]db.JobRef{newRef}, job.JobConsList...),
					}
					database.PutJob(genJob)
					r.Redirect("/a/" + id)
				}
			}
		}
	})

	// About page
	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", nil)
	})

	m.Run()
}
