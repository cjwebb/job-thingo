package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"github.com/microcosm-cc/bluemonday"
	"github.com/satori/go.uuid"

	"github.com/cjwebb/job-thingo/db"
)

type JobForm struct {
	Title        string `form:"Title" binding:"required"`
	Description  string `form:"Description" binding:"required"`
	Rate         string `form:"Rate" binding:"required"`
	ContactEmail string `form:"ContactEmail" binding:"required"`
	UserEmail    string `form:"UserEmail" binding:"required"`
	JobType      string `form:"JobType" binding:"required"`
}

type EmailForm struct {
	Email string `form:"Email" binding:"required"`
}

func main() {
	baseUrl := os.Getenv("BASE_URL")
	database := db.NewDatabase()
	startApp(database, baseUrl)
}

func startApp(database db.Database, baseUrl string) {

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
			sanitizer.Sanitize(form.JobType),
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
	putJob := func(r render.Render, job db.Job) {
		err := database.PutJob(job)
		if err != nil {
			r.HTML(503, "/error", map[string]interface{}{"code": 503})
		} else {
			r.Redirect("/a/"+job.Id+"?display=newLink")
		}
	}

	m.Get("/", func(r render.Render) {
		data := map[string]interface{}{"job": JobForm{}}
		r.HTML(200, "home", data)
	})

	m.Get("/jobs/post", func(r render.Render){
		data := map[string]interface{}{"job": JobForm{}}
		r.HTML(200, "post-job", data)
	})

	m.Post("/jobs/post", binding.Form(JobForm{}), func(form JobForm, r render.Render, errs binding.Errors) {
		saneForm := sanitizeForm(form)
		if errs != nil {
			response := map[string]interface{}{"HasErrors": true, "Errors": errs, "job": saneForm}
			r.HTML(400, "post-job", response)
			return
		}

		id := newId()
		genJob := db.Job{
			Id: id,
			Title: saneForm.Title,
			Description: saneForm.Description,
			ContactEmail: saneForm.ContactEmail,
			Rate: saneForm.Rate,
			JobType: saneForm.JobType,
			JobConsList: []db.JobRef{db.JobRef{Id: id, Email: saneForm.UserEmail}},
		}
		putJob(r, genJob)
	})

	m.Get("/a/:jobid", func(req *http.Request, params martini.Params, r render.Render) {
		jobId, err := uuid.FromString(params["jobid"])
		if err != nil {
			// input not a UUID, so don't try database lookup
			r.Redirect("/")
			return
		}

		job, err := database.GetJob(jobId.String())
		if err != nil {
			log.Print(err.Error())
			r.Redirect("/")
			return
		}

		// if this is a brand new link, then display a panel including the URL
		if req.URL.Query().Get("display") == "newLink" {
			response := map[string]interface{}{"job":job, "displayLink": baseUrl + req.URL.Path}
			r.HTML(200, "job", response)
		} else {
			response := map[string]interface{}{"job":job}
			r.HTML(200, "job", response)
		}
	})

	m.Post("/a/:jobid/gen-link", binding.Form(EmailForm{}), func(form EmailForm, params martini.Params, r render.Render, errs binding.Errors) {
		jobId, err := uuid.FromString(params["jobid"])
		if err != nil {
			// input not a UUID, so don't try database lookup
			r.Redirect("/")
			return
		}

		job, err := database.GetJob(jobId.String())
		if err != nil {
			log.Print(err.Error())
			r.Redirect("/")
			return
		}

		saneForm := sanitizeEmailForm(form)
		if errs != nil {
			response := map[string]interface{}{"HasErrors": true, "Errors": errs, "link": saneForm, "job": job}
			r.HTML(400, "job", response)
			return
		}

		// job exists, form is valid... so gen new link
		id := newId()
		newRef := db.JobRef{Id: id, Email: saneForm.Email}
		genJob := db.Job{
			Id: id,
			Title: job.Title,
			Description: job.Description,
			ContactEmail: job.ContactEmail,
			Rate: job.Rate,
			JobType: job.JobType,
			JobConsList: append([]db.JobRef{newRef}, job.JobConsList...),
		}
		putJob(r, genJob)
	})

	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", nil)
	})

	m.Run()
}
