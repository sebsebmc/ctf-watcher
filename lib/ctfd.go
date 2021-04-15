package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type challenge struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Id       int    `json:"id"`
	PointVal int    `json:"value"`
}

type respWrapper struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

type challengeDetails struct {
	Attempts     int             `json:"attempts"`
	Category     string          `json:"category"`
	Description  string          `json:"description"`
	Files        []string        `json:"files"`
	Hints        []string        `json:"hints"`
	Id           int             `json:"id"`
	Max_attempts int             `json:"max_attempts"`
	Name         string          `json:"name"`
	Solved_by_me bool            `json:"solved_by_me"`
	Solves       int             `json:"solves"`
	State        string          `json:"state"`
	Tags         []string        `json:"tags"`
	Type         string          `json:"type"`
	Type_data    json.RawMessage `json:"type_data"`
	Value        int             `json:"value"`
	View         string          `json:"view"`
}

type CtfdInstance struct {
	url      string
	username string
	password string
	chals    []challenge
	client   http.Client
}

func MakeCtfdInstance(url, username, password string) CtfdInstance {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatalf("Cannot create http client: %v\n", err)
	}

	return CtfdInstance{
		url:    url,
		chals:  nil,
		client: http.Client{Jar: jar},
	}
}

func (ctf *CtfdInstance) GetLatestChallenges() {
	resp, err := ctf.client.Get(ctf.url + "/api/v1/challenges")
	if err != nil {
		log.Fatalf("Failed to get challenges: %v\n", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	wrapper := respWrapper{}
	currentGame := CtfdInstance{}
	latestChals := make([]challenge, 0)
	currentGame.url = ctf.url
	dec.Decode(&wrapper)
	if wrapper.Success {
		json.Unmarshal(wrapper.Data, &latestChals)
		currentGame.chals = latestChals
	}
}

func (ctf *CtfdInstance) getChallengeFiles(chal challenge) challengeDetails {
	resp, err := ctf.client.Get(fmt.Sprintf("%s/api/v1/challenges/%d", ctf.url, chal.Id))
	if err != nil || resp.StatusCode != 200 {
		log.Fatalf("Unable to get chellenge details for challenge %d: %s", chal.Id, chal.Name)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	wrapper := respWrapper{}
	details := challengeDetails{}
	dec.Decode(&wrapper)
	if wrapper.Success {
		json.Unmarshal(wrapper.Data, &details)
		return details
	}
	log.Fatalf("Unable to parse challenge details")
	return challengeDetails{}
}

func (ctf *CtfdInstance) LoginToSite() error {
	values := url.Values{}
	values.Set("name", ctf.username)
	values.Set("password", ctf.password)
	csrfToken, err := ctf.getCSRFToken()
	if err != nil {
		return err
	}
	values.Set("nonce", csrfToken)
	resp, err := ctf.client.PostForm(ctf.url+"/login", values)
	if err != nil {
		return fmt.Errorf("Failed to login: %v\n", err)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to login with status %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

func (ctf *CtfdInstance) getCSRFToken() (string, error) {
	resp, err := ctf.client.Get(ctf.url + "/login")
	if err != nil {
		log.Fatalf("Failed to get login page: %v\n", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to parse document for CSRF token: %v\n", err)
	}
	token, found := doc.Find("input[name=nonce]").First().Attr("value")
	if !found {
		return "", fmt.Errorf("CSRF token not found\n")
	}
	return token, nil
}
