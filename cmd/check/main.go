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
	var request models.ResourceRequest

	// REMOVED HARDCODED

	request.Source.Branch = "develop"
	request.Source.Key = "vh3EaV9qVuEX4Sbk6H"
	request.Source.Secret = "skbEyZ7RLF9YAyZfekjgVxHuJQkAuBce"
	request.Source.URL = "https://api.bitbucket.org"
	request.Source.APIVersion = "2.0"
	request.Source.Team = "lightelligence"
	request.Source.Repo = "notify"
	request.Version.Commit = "15347950458f6b2f1f31202f75cb4b2dda26edce"
	err := logging.PrintText("Unmarshalled struct into", whoami)

	//

	// UNCOMMENT THIS

	//err := json.NewDecoder(os.Stdin).Decode(&request)
	//check(err)
	//err = logging.PrintText("Unmarshalled struct into", whoami)

	//

	check(err)
	err = logging.PrintStruct(request, whoami)
	check(err)

	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)
	commits, err := bitbucket.GetCommitsBranch(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo, request.Source.Branch)
	check(err)

	if request.Version.Commit == "" {
		response = append(response, models.Version{Commit: commits.Values[0].Hash})
	} else {
		for _, commit := range commits.Values {
			if request.Version.Commit == commit.Hash {
				break
			}
			response = append(response, models.Version{Commit: commit.Hash})
		}
	}

	b, _ := json.Marshal(response)
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
