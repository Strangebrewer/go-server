package app

import (
	"github.com/Strangebrewer/go-server/account"
	"github.com/Strangebrewer/go-server/bill"
	"github.com/Strangebrewer/go-server/category"
	"github.com/Strangebrewer/go-server/job"
	"github.com/Strangebrewer/go-server/recruiter"
	"github.com/Strangebrewer/go-server/template"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/transaction"
	"github.com/Strangebrewer/go-server/user"
)

type Application struct {
	ThingStore       *thing.ThingStore
	JobStore         *job.JobStore
	RecruiterStore   *recruiter.RecruiterStore
	TokenStore       *token.TokenStore
	TokenService     *token.TokenService
	UserStore        *user.UserStore
	AccountStore     *account.AccountStore
	CategoryStore    *category.CategoryStore
	BillStore        *bill.BillStore
	TemplateStore    *template.TemplateStore
	TransactionStore *transaction.TransactionStore
}
