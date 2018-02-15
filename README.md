This repo is inspired by [pickledrick/concourse-bitbucket-pullrequest-resource](http://www.github.com/pickledrick/concourse-bitbucket-pullrequest-resource)

A [Concourse](http://concourse.ci/) [resource](http://concourse.ci/resources.html) to interact with the build status API of [Atlassian BitBucket](https://bitbucket.org).

This repo is tied to the [associated Docker image](http://quay.io/pontusarfwedson/concourse-bitbucket-resource) on quay.io, built from the master branch.
## Resource Configuration


These items go in the `source` fields of the resource type. Bold items are required:
 * **`repo`** - repository name to track
 * **`key`** - OAuth key for Consumer
 * **`branch`** - the branch to check and get
 * **`secret`** - OAuth Secret for Consumer
 * **`team`** - Team name repository belongs to
 * **`url`** - bitbucket cloud api path (example: `https://api.bitbucket.org`) **Currently only supported**
 * **`version`** - bitbucket API Version (example: `2.0`) **Currently only supported**
 * **`concourse_url`** - concourse url for setting build link in bitbucket (example: `http://ci.example.com`)



## Behavior


### `check`

Checks for commits on the specified repository and branch. If no version available (first run), will return the HEAD commit. If there is a current version available, will then return the ones newer than the current one.


### `in`

Retrieves a copy of the commit and sets its build state to `INPROGRESS`.

### `out`

Updates the status of a commit.

Parameters:

 * **`commit`** - File containing commit SHA to be updated.
 * **`state`** - the state of the status. Must be one of `success`, `failed` or `inprogress`. By using `inprogress` we can use this resource for updating the status of builds but if you use this resource for get, perform task and then put, then the inprogress updating will be done for you by the `in` script.
 
## Reading logs
The logs can be found by running `fly -t <target-name> intercept -c <pipeline-name>/<resource-name>` and then looking for `check_logfile.txt`, `in_logfile.txt` and `out_logfile.txt` in the `~/root/` directory.

## Example
This resource can be added as a custom resource type, then defined as a resource and then used as follows

```
resource_types:
- name: concourse-bitbucket-resource
  type: docker-image
  source:
    repository: quay.io/pontusarfwedson/concourse-bitbucket-resource

resources:
  - name: ((repo))-test-pr-branch
    type: concourse-bitbucket-resource
    source:
      repo: ((repo))
      branch: "TEST_PR"
      secret: ((bitbucket_secret))
      key: ((bitbucket_key))
      team: teamname
      url: https://api.bitbucket.org
      version: "2.0"
      concourse_url: https://concourse-url.com



jobs:
- name: test-((repo))-poar-bitbucket-resource
  plan:
  - get: ((repo))-test-pr-branch
    trigger: true
  - task: test-((repo))-TEST-PR-branch
    privileged: true
    file: ((repo))-test-pr-branch/build/tasks/task.yaml
    input_mapping: {repo: ((repo))-test-pr-branch} 
    params:
      REPO: ((repo))
      BRANCH: TEST_PR
    on_success:
      put: ((repo))-test-pr-branch
      params:
        commit: ((repo))-test-pr-branch/commit
        state: success
    on_failure:
      put: ((repo))-test-pr-branch
      params:
        commit: ((repo))-test-pr-branch/commit
        state: failed
```


## References

 * [Resources (concourse.ci)](https://concourse.ci/resources.html)
 * [Bitbucket build status API](https://confluence.atlassian.com/bitbucket/use-the-bitbucket-cloud-rest-apis-222724129.html)

## License

[Apache License v2.0]('./LICENSE')
