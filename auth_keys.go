package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type ghKeyEntry struct {
	KeyID   int    `json:"id"`
	KeyData string `json:"key"`
}

func exportAuthorizedKeys(options *RemoteShellOptions) []gossh.PublicKey {

	api_url := fmt.Sprintf("https://api.github.com/users/%s/keys", options.username)

	log.Println()
	log.Println("Exporting Authorized Keys from:", api_url)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.insecureMode,
		},
	}

	githubClient := http.Client{
		Transport: transCfg,
		Timeout:   time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, api_url, nil)
	check(err)

	req.Header.Set("User-Agent", fmt.Sprintf("cloud87-remote-shell/%s", buildVersion))

	res, getErr := githubClient.Do(req)
	check(getErr)

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == 404 {
		log.Fatalf("The user '%s' does not exist on GitHub\n", options.username)
	} else if res.StatusCode != 200 {
		log.Fatalf("Received unexpected status code from Github [%d]\n", res.StatusCode)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	check(readErr)

	var ghkeys []ghKeyEntry

	jsonErr := json.Unmarshal(body, &ghkeys)
	check(jsonErr)

	var keylist []gossh.PublicKey

	for _, element := range ghkeys {

		key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(element.KeyData))
		if err != nil {
			log.Printf("Received error when parsing key. Ignoring. keyid=%d err=%s\n", element.KeyID, err)
		} else {
			fp := gossh.FingerprintSHA256(key)
			log.Println("  Loading Key:", fp)
			keylist = append(keylist, key)
		}
	}

	log.Printf("Loaded %d public keys for user '%s'\n", len(keylist), options.username)

	return keylist

}
