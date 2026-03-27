package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Strangebrewer/go-server/account"
	"github.com/Strangebrewer/go-server/bill"
	"github.com/Strangebrewer/go-server/category"
	"github.com/Strangebrewer/go-server/config"
	"github.com/Strangebrewer/go-server/database"
	"github.com/Strangebrewer/go-server/user"
	"go.mongodb.org/mongo-driver/bson"
)

type SeedData struct {
	Accounts   []account.Account   `json:"accounts"`
	Bills      []bill.Bill         `json:"bills"`
	Categories []category.Category `json:"categories"`
}

func main() {
	cfg := config.InitConfig()
	cfg.LoadEnvVariables()

	mongoConnection, err := database.InitMongoConnection(*cfg)
	if err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}

	seeds, err := os.ReadFile("cmd/seed/seedData.json")
	if err != nil {
		log.Fatalf("osReadFile failed: %v", err)
	}
	var seedData SeedData
	err = json.Unmarshal(seeds, &seedData)
	if err != nil {
		log.Fatalf("unmarshalling failed: %v", err)
	}

	db := mongoConnection.Client.Database(cfg.MongoDBName)

	usersCollection := db.Collection("users")
	var user user.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"email": "derp@test.com"}).Decode(&user)
	if err != nil {
		log.Fatalf("User not found: %v", err)
	}

	accountCollection := db.Collection("accounts")
	accountCollection.DeleteMany(context.TODO(), bson.D{})
	accountsInterface := make([]any, len(seedData.Accounts))
	for i, a := range seedData.Accounts {
		a.UserID = user.ID
		accountsInterface[i] = a
	}
	_, err = accountCollection.InsertMany(context.TODO(), accountsInterface)
	if err != nil {
		log.Fatalf("account insertion failed: %v", err)
	}
	var checkingAccount account.Account
	err = accountCollection.FindOne(context.TODO(), bson.M{"name": "AFCU Checking"}).Decode(&checkingAccount)
	if err != nil {
		log.Fatalf("checking account not found: %v", err)
	}

	billCollection := db.Collection("bills")
	billCollection.DeleteMany(context.TODO(), bson.D{})
	billInterface := make([]any, len(seedData.Bills))
	for i, b := range seedData.Bills {
		b.SourceID = checkingAccount.ID
		b.UserID = user.ID
		billInterface[i] = b
	}
	_, err = billCollection.InsertMany(context.TODO(), billInterface)
	if err != nil {
		log.Fatalf("bill insertion failed: %v", err)
	}

	categoryCollection := db.Collection("categories")
	categoryCollection.DeleteMany(context.TODO(), bson.D{})
	categoryInterface := make([]any, len(seedData.Categories))
	for i, c := range seedData.Categories {
		c.UserID = user.ID
		categoryInterface[i] = c
	}
	_, err = categoryCollection.InsertMany(context.TODO(), categoryInterface)
	if err != nil {
		log.Fatalf("category insertion failed: %v", err)
	}
}
