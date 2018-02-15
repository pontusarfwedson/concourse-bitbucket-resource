package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/models"
)

// GetCommitsBranch gets all the commits for a specific repo and branch
func GetCommitsBranch(url string, token string, version string, team string, repo string, branch string) (*models.CommitsResponse, error) {
	req, err := http.NewRequest("GET", url+"/"+version+"/repositories/"+team+"/"+repo+"/commits/"+branch, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create http request")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	var client = &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Could not run http get")
	}
	defer response.Body.Close()

	var cmtRsp models.CommitsResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&cmtRsp)

	return &cmtRsp, err
}

// SetBuildStatus updates the commit associated with a pull-request and sets the state () as well as a link to the Concourse build log.
func SetBuildStatus(url, token, version, team, repo, commit, state, concourseHost string) error {
	if url == "" {
		return errors.New("url must be provided")
	}
	if token == "" {
		return errors.New("token must be provided")
	}
	if version == "" {
		return errors.New("version must be provided")
	}
	if team == "" {
		return errors.New("team must be provided")
	}
	if repo == "" {
		return errors.New("repo must be provided")
	}
	if commit == "" {
		return errors.New("commit must be provided")
	}
	if state == "" {
		return errors.New("state must be provided")
	}
	if concourseHost == "" {
		return errors.New("concourse host must be provided")
	}

	buildJob := os.Getenv("BUILD_JOB_NAME")

	concourseURL := fmt.Sprintf(
		"%s/teams/%s/pipelines/%s/jobs/%s/builds/%s",
		concourseHost,
		os.Getenv("BUILD_TEAM_NAME"),
		os.Getenv("BUILD_PIPELINE_NAME"),
		buildJob,
		os.Getenv("BUILD_NAME"),
	)

	key := "concourse-" + buildJob
	if len(key) >= 43 {
		key = key[0:40]
	}

	status := models.OutStatus{State: state, Key: key, URL: concourseURL}
	out, err := json.Marshal(status)
	if err != nil {
		return errors.Wrapf(err, "unable to marshal build status: %+v", status)
	}

	req, err := http.NewRequest("POST", url+"/"+version+"/repositories/"+team+"/"+repo+"/commit/"+commit+"/statuses/build", bytes.NewBuffer(out))
	if err != nil {
		return errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request to set build status failed")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			return errors.Wrapf(err, "request to set build status failed with status [%d], but the response body could not be read", res.StatusCode)
		}
		return errors.Errorf("request to set build status failed, code [%d], url [%s], body: %s", res.StatusCode, req.URL, buf.String())
	}
	return nil
}

func RequestToken(key string, secret string) (string, error) {
	if key == "" {
		return "", errors.New("key must be provided")
	}
	if secret == "" {
		return "", errors.New("secret must be provided")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "unable to create request object")
	}
	req.SetBasicAuth(key, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var response models.Token
	err = do(req, &response)
	if err != nil {
		return "", errors.Wrap(err, "request for token failed")
	}
	return response.AccessToken, nil
}
