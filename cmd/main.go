package main

import (
	"log"

	"github.com/Chrisentech/aluta-market-api/db"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/Chrisentech/aluta-market-api/router"
)

func main() {
	dbConn, err := db.NewDatabse()
	if err != nil {
		log.Fatalf("could not initialize database connection %s \n", err)
	}
	userRep := user.NewRespository(dbConn.GetDB())
	userSrvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSrvc)
	router.InitRouter(userHandler)
	router.Start("0.0.0.0:8082")
}

// https://www.youtube.com/watch?v=cphghdh1DoY&list=PLzQWIQOqeUSNwXcneWYJHUREAIucJ5UZn&index=2
