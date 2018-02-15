package models

import "time"

type CommitsResponse struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
		Hash       string `json:"hash"`
		Repository struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Comments struct {
				Href string `json:"href"`
			} `json:"comments"`
			Patch struct {
				Href string `json:"href"`
			} `json:"patch"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Diff struct {
				Href string `json:"href"`
			} `json:"diff"`
			Approve struct {
				Href string `json:"href"`
			} `json:"approve"`
			Statuses struct {
				Href string `json:"href"`
			} `json:"statuses"`
		} `json:"links"`
		Author struct {
			Raw  string `json:"raw"`
			Type string `json:"type"`
			User struct {
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
				Type        string `json:"type"`
				UUID        string `json:"uuid"`
				Links       struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
					Avatar struct {
						Href string `json:"href"`
					} `json:"avatar"`
				} `json:"links"`
			} `json:"user"`
		} `json:"author"`
		Parents []struct {
			Hash  string `json:"hash"`
			Type  string `json:"type"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
		} `json:"parents"`
		Date    time.Time `json:"date"`
		Message string    `json:"message"`
		Type    string    `json:"type"`
	} `json:"values"`
	Next string `json:"next"`
}

// Links is the structure of links and references attached to many Bitbucket API responses.
type Links struct {
	Activity struct {
		Href string `json:"href"`
	} `json:"activity"`
	Approve struct {
		Href string `json:"href"`
	} `json:"approve"`
	Avatar struct {
		Href string `json:"href"`
	} `json:"avatar"`
	Comments struct {
		Href string `json:"href"`
	} `json:"comments"`
	Commits struct {
		Href string `json:"href"`
	} `json:"commits"`
	Decline struct {
		Href string `json:"href"`
	} `json:"decline"`
	Diff struct {
		Href string `json:"href"`
	} `json:"diff"`
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
	Merge struct {
		Href string `json:"href"`
	} `json:"merge"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Statuses struct {
		Href string `json:"href"`
	} `json:"statuses"`
}

// CheckRequest is the struct/JSON that is supplied to "check", coming from the Concourse pipeline under "resources"
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// CheckResponse is the struct/JSON returned from "check"
type CheckResponse []Version

// Params ... (referenced from OutRequest)
type Params struct {
	State       string `json:"state"`
	PullRequest string `json:"pull_request"`
	Commit      string `json:"commit"`
}

// Source ... (referenced from CheckRequest)
type Source struct {
	Branch       string `json:"branch"`
	Repo         string `json:"repo"`
	Secret       string `json:"secret"`
	Key          string `json:"key"`
	Team         string `json:"team"`
	ConcourseURL string `json:"concourse_url"`
}

// Version ... (referenced from CheckRequest)
type Version struct {
	Commit string `json:"commit"`
}

// InRequest is the struct/JSON supplied as input to "in" - Concourse pipeline "get"
type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// InResponse is the struct/JSON that is output from "in".
type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

// Metadata holds multiple MetadataField, which is output from "in" and "out"
type Metadata []MetadataField

// MetadataField holding data presented as metedata in "in" and "out"
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OutRequest is the struct/JSON supplied as input to "out" - Concourse pipeline "put"
type OutRequest struct {
	Params  Params  `json:"params"`
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// OutResponse is the struct/JSON that is output from "out".
type OutResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

// OutStatus holds data about a build's status.
type OutStatus struct {
	State string `json:"state"`
	Key   string `json:"key"`
	URL   string `json:"url"`
}

// type CredentialsRequest2 struct {
// 	GrantType string `json:"grant_type"`
// }

// Token holds Authentication Tokens for accessing the Bitbucket API.
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scopes       string `json:"scopes"`
	TokenType    string `json:"token_type"`
}

// CommitResponse represents the Commit Status response from the Bitbucket API.
// <https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/commit/%7Bnode%7D/statuses>
type CommitStatusResponse struct {
	Page    int `json:"page"`
	Pagelen int `json:"pagelen"`
	Size    int `json:"size"`
	Values  []struct {
		CreatedOn   time.Time   `json:"created_on"`
		Description string      `json:"description"`
		Key         string      `json:"key"`
		Links       Links       `json:"links"`
		Name        string      `json:"name"`
		Refname     interface{} `json:"refname"`
		Repository  struct {
			FullName string `json:"full_name"`
			Links    struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Name string `json:"name"`
			Type string `json:"type"`
			UUID string `json:"uuid"`
		} `json:"repository"`
		State     string `json:"state"`
		Type      string `json:"type"`
		UpdatedOn string `json:"updated_on"`
		URL       string `json:"url"`
	} `json:"values"`
}
