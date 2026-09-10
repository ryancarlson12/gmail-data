package main

// NOTE: gmail access tokens refresh every 1hr.

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func query(accessToken string, url string) error {
	// get the curl command
	cmd := exec.Command("curl", "-s", "-H", fmt.Sprintf("Authorization: Bearer %s", accessToken), url)

	// Get output and errs
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ERROR: curl execution failed: %w (stderr: %s)", err, stderr.String())
		// create an error message here specifically for when the token expired
	}

	fmt.Println(stdout.String())
	return nil
}

func main() {
	// valid command line flags
	actionFlag := flag.String("action", "profile", "Action to perform (supports: profile)")

	flag.Parse()

	// must set GMAIL_ACCESS_TOKEN environment variable
	accessToken := os.Getenv("GMAIL_ACCESS_TOKEN")
	fmt.Println(accessToken)

	if accessToken == "" {
		fmt.Println("ERROR: GMAIL_ACCESS_TOKEN not defined. Try 'export GMAIL_ACCESS_TOKEN=\"<your access token>\"")
		os.Exit(1)
	}

	switch *actionFlag {
	case "profile":
		fmt.Println("Fetching Gmail profile...")
		if err := query(accessToken, "https://gmail.googleapis.com/gmail/v1/users/me/profile"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: could not fetch Gmail profile: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("ERROR: unknown action: %v\n", *actionFlag)
		os.Exit(1)
	}
}
