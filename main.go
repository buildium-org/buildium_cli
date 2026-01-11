package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	case "tutorial":
		tutorial(os.Args[2:])
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

func tutorial(args []string) {
	ctx := context.Background()
	client := supabase.NewSupaClient(ctx)
	err := client.VerifyAuthToken(ctx)
	if err != nil {
		fmt.Printf("Failed to verify auth token: %v\n", err)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Println("Usage: buildium tutorial <command> [options]")
		fmt.Println("Commands: create-template")
		os.Exit(1)
	}

	switch args[0] {
	case "create-template":
		tutorialTemplate(args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		os.Exit(1)
	}
}

func tutorialTemplate(args []string) {
	tutorialFlags := flag.NewFlagSet("tutorial create-template", flag.ExitOnError)
	repoName := tutorialFlags.String("name", "", "Repository name for the template")
	tutorialFlags.Parse(args)

	if *repoName == "" {
		fmt.Println("Repository name is required.")
		fmt.Println("Usage: buildium tutorial create-template --name <repository name>")
		os.Exit(1)
	}

	fmt.Printf("Creating starter template for repository name: %s\n", *repoName)

	githubUrl := "https://github.com/buildium-org/tutorial_template.git"

	cmd := exec.Command("git", "clone", githubUrl, *repoName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to clone repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully cloned repository.")

	// Go through all files in the cloned repo and replace <YOUR_IMAGE_NAME_HERE> with the real image name
	err := replaceInDirectory(*repoName, "<YOUR_IMAGE_NAME_HERE>", *repoName)
	if err != nil {
		fmt.Printf("Failed to update <YOUR_IMAGE_NAME_HERE>: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated <YOUR_IMAGE_NAME_HERE>.")
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
	repoName := projectFlags.String("name", "", "Repository name for the template")
	lang := projectFlags.String("lang", "", "Programming language for the template")
	projectFlags.Parse(args)

	if *projectID == "" || *lang == "" || *repoName == "" {
		fmt.Println("Project ID, language, and repository name are required.")
		fmt.Println("Usage: buildium project create-template --projectid <id> --lang <language> --name <repository name>")
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

	fmt.Printf("Project ID: %s, Project Name: %s, Tutorial Docker Image: %s, language: %s\n", project.ProjectId, *repoName, tutorial.DockerImage, *lang)

	var githubUrl string
	switch *lang {
	case "go":
		fmt.Println("Creating Go starter template")
		githubUrl = "https://github.com/buildium-org/go_template.git"
	case "typescript":
		fmt.Println("Creating TypeScript starter template")
		githubUrl = "https://github.com/buildium-org/ts_template.git"
	default:
		fmt.Printf("Unknown language: %s\n", *lang)
		os.Exit(1)
	}

	cmd := exec.Command("git", "clone", githubUrl, *repoName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to clone repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully cloned repository.")

	// Go through all files in the cloned repo and replace <YOUR_PROJECT_ID> with the real projectID

	err = replaceInDirectory(*repoName, "<YOUR_PROJECT_ID>", project.ProjectId)
	if err != nil {
		fmt.Printf("Failed to update <YOUR_PROJECT_ID>: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated <YOUR_PROJECT_ID>.")
	err = replaceInDirectory(*repoName, "<YOUR_IMAGE_NAME_HERE>", *repoName)
	if err != nil {
		fmt.Printf("Failed to update <YOUR_IMAGE_NAME_HERE>: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated <YOUR_IMAGE_NAME_HERE>.")

	err = replaceInDirectory(*repoName, "<TEST_HARNESS_IMAGE_HERE>", tutorial.DockerImage)
	if err != nil {
		fmt.Printf("Failed to update <TEST_HARNESS_IMAGE_HERE>: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully updated <TEST_HARNESS_IMAGE_HERE>.")
}

func replaceInDirectory(repoPath string, old string, new string) error {
	return filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only consider regular files
		if info.Mode().IsRegular() {
			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(input)
			changed := strings.ReplaceAll(content, old, new)
			if content != changed {
				err = os.WriteFile(path, []byte(changed), info.Mode().Perm())
				if err != nil {
					return err
				}
				fmt.Printf("Updated %s in: %s\n", old, path)
			}
		}
		return nil
	})
}
