package lib

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Creds struct {
	Url      string
	Username string
	Password string
}

func GetCreds() (*Creds, error) {
	contents, err := ioutil.ReadFile("./.ctf")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("Creds file is corrupted")
	}
	return &Creds{
		Url:      lines[0],
		Username: lines[1],
		Password: lines[2],
	}, nil
}
