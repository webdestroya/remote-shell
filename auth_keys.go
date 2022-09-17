package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

const sshAuthKeyEnvVar = "C87_RSHELL_AUTHORIZED_KEY"

type ghKeyEntry struct {
	KeyID   int    `json:"id"`
	KeyData string `json:"key"`
}

func isAuthKeyFromEnv() bool {
	return os.Getenv(sshAuthKeyEnvVar) != ""
}

func exportAuthorizedKeys() []gossh.PublicKey {
	var keylist []gossh.PublicKey

	if isAuthKeyFromEnv() {
		key := exportAuthorizedKeysFromEnv()
		if key != nil {
			keylist = append(keylist, key)
		}
	}

	if globalOptions.username != "" {
		ghkeylist := exportAuthorizedKeysFromGitHub()
		keylist = append(keylist, ghkeylist...)
	}

	return keylist
}

func exportAuthorizedKeysFromEnv() gossh.PublicKey {

	log.Println("SSH Key provided via environment variable")

	envKeyValue := os.Getenv(sshAuthKeyEnvVar)

	key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(envKeyValue))
	if err != nil {
		log.Printf("Received error when parsing key. Ignoring. err=%s\n", err)
	} else {
		fp := gossh.FingerprintSHA256(key)
		log.Println("Loading Key:", fp)
		return key
	}
	return nil
}

func exportAuthorizedKeysFromGitHub() []gossh.PublicKey {

	ghKeyListUrl := fmt.Sprintf("https://api.github.com/users/%s/keys", globalOptions.username)

	log.Println()
	log.Println("Exporting Authorized Keys from:", ghKeyListUrl)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: globalOptions.insecureMode,
		},
	}

	githubClient := http.Client{
		Transport: transCfg,
		Timeout:   time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, ghKeyListUrl, nil)
	check(err)

	req.Header.Set("User-Agent", fmt.Sprintf("cloud87-remote-shell/%s", buildVersion))

	res, err := githubClient.Do(req)
	check(err)

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == 404 {
		panic(fmt.Sprintf("The user '%s' does not exist on GitHub", globalOptions.username))
	} else if res.StatusCode != 200 {
		panic(fmt.Sprintf("Received unexpected status code from Github [%d]", res.StatusCode))
	}

	var ghkeys []ghKeyEntry

	err = json.NewDecoder(res.Body).Decode(&ghkeys)
	check(err)

	var keylist []gossh.PublicKey

	for _, element := range ghkeys {

		key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(element.KeyData))
		if err != nil {
			log.Printf("Received error when parsing key. Ignoring. keyid=%d err=%s\n", element.KeyID, err)
		} else {
			fp := gossh.FingerprintSHA256(key)
			log.Println("Loading Key:", fp)
			keylist = append(keylist, key)
		}
	}

	log.Printf("Loaded %d public keys for user '%s'\n", len(keylist), globalOptions.username)

	return keylist

}
