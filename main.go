package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"buildium_cli/supabase"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: buildium <command> [options]")
		fmt.Println("Commands: login")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		login(os.Args[2:])
	case "project":
		project(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func login(args []string) {
	ctx := context.Background()
	client := supabase.NewSupaClient(ctx)
	loginFlags := flag.NewFlagSet("login", flag.ExitOnError)
	email := loginFlags.String("email", "", "email to use for login")
	password := loginFlags.String("password", "", "password to use for login")
	loginFlags.Parse(args)

	if *email == "" || *password == "" {
		fmt.Println("Email and password are required `buildium login -email <email> -password <password>`")
		os.Exit(1)
	}

	fmt.Printf("Logging in with email: %s\n", *email)

	err := client.Login(ctx, *email, *password)
	if err != nil {
		fmt.Printf("Failed to login: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Logged in successfully, token: ", client.Config.AuthToken)
}

func project(args []string) {
	ctx := context.Background()
	client := supabase.NewSupaClient(ctx)
	err := client.VerifyAuthToken(ctx)
	if err != nil {
		fmt.Printf("Failed to verify auth token: %v\n", err)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Println("Usage: buildium project <command> [options]")
		fmt.Println("Commands: list")
		os.Exit(1)
	}

	switch args[0] {
	case "create-template":
		projectTemplate(args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}

func projectTemplate(args []string) {
	projectFlags := flag.NewFlagSet("project create-template", flag.ExitOnError)
	projectID := projectFlags.String("projectid", "", "Project ID to use for the template")
	lang := projectFlags.String("lang", "", "Programming language for the template")
	projectFlags.Parse(args)

	if *projectID == "" || *lang == "" {
		fmt.Println("Project ID and language are required.")
		fmt.Println("Usage: buildium project create-template --projectid <id> --lang <language>")
		os.Exit(1)
	}

	switch *lang {
	case "go":
		fmt.Println("Creating Go starter template")
	case "typescript":
		fmt.Println("Creating TypeScript starter template")
	default:
		fmt.Printf("Unknown language: %s\n", *lang)
		os.Exit(1)
	}

	fmt.Printf("Creating starter template for project ID: %s and language: %s\n", *projectID, *lang)

	ctx := context.Background()
	client := supabase.NewSupaClient(ctx)
	fmt.Println("Getting project: ", *projectID)
	project, err := client.GetProject(ctx, *projectID)
	if err != nil {
		fmt.Printf("Failed to get project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully got project: ", project.ProjectId)
	fmt.Println("Getting tutorial: ", project.TutorialId)
	tutorial, err := client.GetTutorial(ctx, project.TutorialId)
	if err != nil {
		fmt.Printf("Failed to get tutorial: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully got tutorial: ", tutorial.TutorialId)

	fmt.Printf("Project ID: %s, Project Image Name: %s, Tutorial Docker Image: %s, language: %s\n", project.ProjectId, project.Name, tutorial.DockerImage, *lang)

	// TODO: get template based on language and then do a find and replace across the template for these retrieved values
}
