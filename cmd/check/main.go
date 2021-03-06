package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/bitbucket"
	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/logging"
	"github.com/pontusarfwedson/concourse-bitbucket-resource/cmd/models"
)

const (
	whoami logging.ResourceModule = logging.Check
)

func main() {
	var response models.CheckResponse
	var request models.InRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	check(err)

	err = logging.PrintText("Unmarshalled struct into", whoami)
	check(err)

	err = logging.PrintStruct(request, whoami)
	check(err)

	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)
	commits, err := bitbucket.GetCommitsBranch(bitbucket.Url, token, bitbucket.ApiVersion, request.Source.Team, request.Source.Repo, request.Source.Branch)
	check(err)

	if request.Version.Commit == "" && len(commits.Values) > 0 {
		response = append(response, models.Version{Commit: commits.Values[0].Hash})
	} else {
		for _, commit := range commits.Values {
			if request.Version.Commit == commit.Hash {
				response = append(response, models.Version{Commit: commit.Hash})
				break
			}
			response = append(response, models.Version{Commit: commit.Hash})
		}
	}

	//Revert our response so that we get the newer versions below the older ones
	responseReverted := make(models.CheckResponse, len(response), len(response))
	index := 0
	for i := len(response) - 1; i >= 0; i-- {
		responseReverted[i] = response[index]
		index++
	}

	b, _ := json.Marshal(responseReverted)
	jsonStr := string(b)
	err = logging.PrintText(fmt.Sprintf(">>>>>>>>>>     Output to os.Stdout is %s", jsonStr), whoami)
	check(err)
	fmt.Fprintf(os.Stdout, jsonStr)

}

func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
