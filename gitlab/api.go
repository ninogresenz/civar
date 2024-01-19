package gitlab

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
)

// Api provides necessary methods to interact with the Gitlab's REST API.
// Docs: https://docs.gitlab.com/ee/api/api_resources.html
type Api interface {
	Search(term string) ([]Project, error)
	GetProjectVars(project string) (CiVariableList, error)
	CreateVar(project string, variable CiVariable) (*CiVariable, error)
	UpdateVar(project string, variable CiVariable) (*CiVariable, error)
}

type api struct {
	gitlabUrl string
	api       *sling.Sling
}

func New(gitlabUrl string, token string, httpClient *http.Client) Api {
	return api{
		gitlabUrl: gitlabUrl,
		api:       sling.New().Base(gitlabUrl).Client(httpClient).Set("PRIVATE-TOKEN", token),
	}
}

func (a api) Search(term string) (projects []Project, err error) {
	req := a.api.New().Get(fmt.Sprintf("/api/v4/search?scope=projects&search=%s", term))
	return paginate[Project](req, projects)
}

func (a api) GetProjectVars(project string) (allVars CiVariableList, err error) {
	projectEncoded := url.QueryEscape(project)
	req := a.api.New().Get(fmt.Sprintf("/api/v4/projects/%s/variables", projectEncoded))
	return paginate[CiVariable](req, allVars)
}

func (a api) CreateVar(project string, variable CiVariable) (created *CiVariable, err error) {
	projectEncoded := url.QueryEscape(project)
	var errorResponse ErrorResponse
	resp, err := a.api.New().
		Post(fmt.Sprintf("/api/v4/projects/%v/variables", projectEncoded)).
		BodyJSON(variable).
		Receive(&created, &errorResponse)
	err = handleHttpError(resp, err, errorResponse)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (a api) UpdateVar(project string, variable CiVariable) (updated *CiVariable, err error) {
	projectEncoded := url.QueryEscape(project)
	var errorResponse ErrorResponse
	body := UpdateBody{variable.Value, Filter{variable.EnvironmentScope}}
	resp, err := a.api.New().
		Put(fmt.Sprintf("/api/v4/projects/%v/variables/%v", projectEncoded, variable.Key)).
		BodyForm(&body).
		Receive(&updated, &errorResponse)
	err = handleHttpError(resp, err, errorResponse)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
