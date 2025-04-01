package gerritapi

import (
	"context"

	"github.com/andygrunwald/go-gerrit"
)

func InitGerritClient(url string, user string, password string) *gerrit.Client {

	gerritClient, err := gerrit.NewClient(context.Background(), url, nil)

	if err != nil {
		panic("failed to init gerrit client.")
	}
	gerritClient.Authentication.SetBasicAuth(user, password)
	return gerritClient
}

func GetChangeInclueCurrentRevison(gerritClient *gerrit.Client, changeID string) (changeInfo *gerrit.ChangeInfo, error error) {
	opt := &gerrit.ChangeOptions{}
	opt.AdditionalFields = []string{"CURRENT_REVISION"}

	changeInfo, _, err := gerritClient.Changes.GetChangeDetail(context.Background(), changeID, opt)
	if err != nil {
		return nil, err
	}
	return changeInfo, nil
}

func GetFileListByCommit(gerritClient *gerrit.Client, changeID string, revision string) ([]string, error) {

	filesInfo, _, err := gerritClient.Changes.ListFiles(context.Background(), changeID, revision, nil)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for key := range filesInfo {
		files = append(files, key)
	}

	return files, nil

}

func GetDiff(gerritClient *gerrit.Client, changeID string, revision string, fileId string) (*gerrit.DiffInfo, error) {
	opt := &gerrit.DiffOptions{}
	diffInfo, _, err := gerritClient.Changes.GetDiff(context.Background(), changeID, revision, fileId, opt)
	if err != nil {
		return nil, err
	}
	return diffInfo, nil
}

func GetChildProject(gerritClient *gerrit.Client, projectName string) ([]string, error) {
	opt := &gerrit.ChildProjectOptions{}
	ChildProject, _, err := gerritClient.Projects.ListChildProjects(context.Background(), projectName, opt)
	if err != nil {
		return nil, err
	}
	project := []string{}
	for _, value := range *ChildProject {
		project = append(project, value.Name)
	}
	return project, nil
}

type Proj struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
}

func GetProjectsByBranch(gerritClient *gerrit.Client, branch string) ([]Proj, error) {

	opt := &gerrit.ProjectOptions{}
	opt.Branch = branch

	projectList, _, err := gerritClient.Projects.ListProjects(context.Background(), opt)
	if err != nil {
		return nil, err
	}

	projects := make([]Proj, 0)
	for key, value := range *projectList {
		projects = append(projects, Proj{Name: key, Revision: value.Branches[branch]})
	}
	return projects, nil
}

func GetParentProject(gerritClient *gerrit.Client, name string) (string, error) {
	project, _, err := gerritClient.Projects.GetProjectParent(context.Background(), name)
	if err != nil {
		return "", err
	}
	return project, nil

}
