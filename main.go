package main

import (
	"flag"
	"fmt"
	"log"
)

func seedUser(store Storage, fname, lname, email, pw string) *User {
	user, err := NewUser(fname, lname, email, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateUser(user); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new user=> ", user.FirstName, user.LastName)

	return user
}

func Seed(s Storage) {
	seedUser(s, "umut", "okur", "example@gmail.com", "passwd")
}

func main() {
	seed := flag.Bool("seed", true, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		Seed(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
