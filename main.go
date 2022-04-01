package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var (
	file     = flag.String("file", "", "")
	domain   = flag.String("domain", "", "")
	username = flag.String("username", "", "")
	password = flag.String("password", "", "")
)

func main() {

	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	token := getToken()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addRecipe(token, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getToken() token {

	data := url.Values{}
	data.Set("grant_type", "")
	data.Set("scope", "")
	data.Set("client_id", "")
	data.Set("client_secret", "")
	data.Set("username", *username)
	data.Set("password", *password)

	u, err := url.Parse(*domain)
	u.Path = path.Join(u.Path, "/api/auth/token")

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatalln(err)
	}

	return token
}

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func addRecipe(token token, recipe string) {

	if !strings.HasPrefix(recipe, "http") {
		log.Println("Invalid URL: ", recipe)
		return
	}

	b, err := json.Marshal(map[string]string{"url": recipe})
	if err != nil {
		log.Fatalln(err)
	}

	u, err := url.Parse(*domain)
	u.Path = path.Join(u.Path, "/api/recipes/create-url")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("authorization", "Bearer "+token.AccessToken)
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}
