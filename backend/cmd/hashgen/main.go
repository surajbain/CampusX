package main

import (
	"fmt"
	"os"

	"github.com/campusx/api/pkg/crypto"
)

func main() {
	passwords := map[string]string{
		"super@campusx.dev":     "SuperAdmin@123",
		"admin@iitb.edu":        "Admin@123",
		"organizer@iitb.edu":    "Organizer@123",
		"volunteer@iitb.edu":    "Volunteer@123",
		"student@iitb.edu":      "Student@123",
		"advertiser@brandco.in": "Advertiser@123",
	}

	for email, pw := range passwords {
		hash, err := crypto.Hash(pw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash %s: %v\n", email, err)
			os.Exit(1)
		}
		fmt.Printf("%s -> %s\n", email, hash)
	}
}
