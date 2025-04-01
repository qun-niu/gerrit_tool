package main

import (
	"encoding/json"
	"fmt"
	"gerrrit_tool/gerritapi"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var (
	project           = kingpin.Command("project", "Lists the projects accessible by the caller.")
	branch            = project.Flag("with_branch", "Get projects that have a special branch").String()
	childProject      = kingpin.Command("childproject", "Lists the child project of project")
	parentProjectName = childProject.Arg("name", "parent project name").Required().String()
	parentProject     = kingpin.Command("parentproject", "Retrieves the name of a project's parent project.")
	childProjectName  = parentProject.Arg("name", "child project name").Required().String()
)

func main() {

	url := os.Getenv("TOOL_GERRIT_URL")
	user := os.Getenv("TOOL_GERRIT_USER")
	password := os.Getenv("TOOL_GERRIT_PASSWORD")

	if url == "" || user == "" || password == "" {
		log.Fatal("please set 'TOOL_GERRIT_URL', 'TOOL_GERRIT_USER', 'TOOL_GERRIT_PASSWORD'")
	}

	gerritClient := gerritapi.InitGerritClient(url, user, password)
	switch kingpin.Parse() {
	case "project":
		if *branch != "" {
			projectList, err := gerritapi.GetProjectsByBranch(gerritClient, *branch)
			data, _ := json.Marshal(projectList)
			if err != nil {
				log.Fatal(fmt.Printf("failed to get project: %s", err.Error()))
			}
			fmt.Println(string(data))
		}

	case "childproject":
		projectList, err := gerritapi.GetChildProject(gerritClient, *parentProjectName)
		if err != nil {
			log.Fatal(fmt.Printf("failed to get child projects: [parent: %s]  %s", *parentProjectName, err.Error()))
		}
		fmt.Println(projectList)

	case "parentproject":

		if *childProjectName != "" {
			project, err := gerritapi.GetParentProject(gerritClient, *childProjectName)
			if err != nil {
				log.Fatal()
			}
			fmt.Println(project)
		}

	}

}
