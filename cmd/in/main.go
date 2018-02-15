package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/bitbucket"
	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/logging"
	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/models"
)

const (
	whoami logging.ResourceModule = logging.In
)

func main() {

	var request models.InRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	check(err)

	if request.Version.Commit == "" {
		log.Printf("Ignoring input request without version (commit)")
		err = json.NewEncoder(os.Stdout).Encode(models.InResponse{Version: models.Version{}, Metadata: models.Metadata{}})
		check(err)
		return
	}
	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)

	err = bitbucket.SetBuildStatus(bitbucket.Url, token, bitbucket.ApiVersion, request.Source.Team, request.Source.Repo, request.Version.Commit, "INPROGRESS", request.Source.ConcourseURL)
	check(err)

	args := os.Args

	outputDir := args[1]
	commitID := []byte(string(strings.Replace(request.Version.Commit, "\n", "", -1)))

	err = os.MkdirAll(outputDir, os.ModePerm)
	check(err)

	r, err := git.PlainClone(outputDir, false, &git.CloneOptions{
		URL: "https://x-token-auth:" + token + "@bitbucket.org/" + request.Source.Team + "/" + request.Source.Repo,
	})
	check(err)

	w, err := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{

		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", request.Source.Branch)),
		Force:  true,
	})

	err = ioutil.WriteFile(outputDir+"/commit", commitID, 0644)
	check(err)

	version := models.MetadataField{Name: "Commit", Value: request.Version.Commit}
	metadata := models.Metadata{version}

	err = json.NewEncoder(os.Stdout).Encode(models.InResponse{Version: request.Version, Metadata: metadata})
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
