package server

import (
	"context"

	"github.com/Strangebrewer/go-server/account"
	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/bill"
	"github.com/Strangebrewer/go-server/category"
	"github.com/Strangebrewer/go-server/job"
	"github.com/Strangebrewer/go-server/recruiter"
	tmpl "github.com/Strangebrewer/go-server/template"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/transaction"
	"github.com/Strangebrewer/go-server/user"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	// Ensure TTL index exists (do once at startup)
	_ = application.TokenStore.EnsureIndexes(context.Background())
	_ = application.UserStore.EnsureIndexes(context.Background())

	r.Mount("/things", thing.ThingRoutes(application.ThingStore, application.TokenService))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/jobs", job.JobRoutes(application.JobStore, application.RecruiterStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/recruiters", recruiter.RecruiterRoutes(application.RecruiterStore))
	r.Mount("/token", token.TokenRoutes(application.TokenService))
	r.Mount("/users", user.UserRoutes(application.UserStore, application.TokenService))

	r.With(token.RequireAccess(application.TokenService)).
		Mount("/accounts", account.AccountRoutes(application.AccountStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/categories", category.CategoryRoutes(application.CategoryStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/bills", bill.BillRoutes(application.BillStore, application.TransactionStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/templates", tmpl.TemplateRoutes(application.TemplateStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/transactions", transaction.TransactionRoutes(application.TransactionStore))
	// r.With(token.RequireAccess(application.TokenService)).
	// 	Mount("/dashboard", dashboard.DashboardRoutes(application.MemberStore, application.TransactionStore))
}
