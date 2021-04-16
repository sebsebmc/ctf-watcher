package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type ctfdChallenge struct {
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
	csrf     string
	chals    []ctfdChallenge
	client   http.Client
}

// MakeCtfdInstance creates an object that can be used to interact with a
// CTFd instance. You must call LoginToSite on this before issuing other calls
// TODO: check login state via session cookie expiry and auto-attempt login
func MakeCtfdInstance(url, username, password string) CtfdInstance {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatalf("Cannot create http client: %v\n", err)
	}

	return CtfdInstance{
		url:      url,
		username: username,
		password: password,
		chals:    nil,
		client:   http.Client{Jar: jar},
	}
}

func (ctf *CtfdInstance) GetLatestChallenges() error {
	req, _ := http.NewRequest("GET", ctf.url+"/api/v1/challenges", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("csrf-token", ctf.csrf)
	req.Header.Set("User-Agent", "ctf-watcher/dev")
	resp, err := ctf.client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get challenges: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to get challenges: server responded with %d\n", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	wrapper := respWrapper{}
	latestChals := make([]ctfdChallenge, 0)
	err = dec.Decode(&wrapper)
	if err != nil {
		return err
	}
	if wrapper.Success {
		json.Unmarshal(wrapper.Data, &latestChals)
		ctf.chals = latestChals
	} else {
		return fmt.Errorf("API call to api/v1/challenges failed")
	}
	return nil
}

func (ctf *CtfdInstance) PrintChallenges() {
	for _, v := range ctf.chals {
		fmt.Println(v.Name)
	}
}

func (ctf *CtfdInstance) getChallengeFiles(chal ctfdChallenge) challengeDetails {
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
	loginToken, err := ctf.getLoginNonce()
	if err != nil {
		return err
	}
	ctf.csrf = loginToken
	values.Set("nonce", loginToken)
	resp, err := ctf.client.PostForm(ctf.url+"/login", values)
	if err != nil {
		return fmt.Errorf("Failed to login: %v\n", err)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to login with status %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	csrfLocator := []byte("'csrfNonce': \"")
	body, _ := ioutil.ReadAll(resp.Body)
	start := bytes.Index(body, csrfLocator)
	csrfFound := false
	if start != -1 {
		body = body[start+len(csrfLocator):]
		end := bytes.Index(body, []byte{'"'})
		if end != -1 {
			ctf.csrf = string(body[:end])
			csrfFound = true
		}
	}
	if !csrfFound {
		return errors.New("Unable to find csrfToken")
	}
	return nil
}

func (ctf *CtfdInstance) getLoginNonce() (string, error) {
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
