package bitbucket

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

// do will perform a http request with retries and backoff
// will then unmarshall into the passed response object
func do(request *http.Request, response interface{}) error {
	client := pester.New()
	client.MaxRetries = 10
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, client.LogString())
	}

	// Some Bitbucket APIs can return 404 in some cases.
	if resp.StatusCode == 404 {
		return nil
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.Errorf("request failed, code [%d], url [%s], body: %s", resp.StatusCode, request.URL, buf.String())
	}
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal the response: %s", buf.String())
	}
	return nil
}
