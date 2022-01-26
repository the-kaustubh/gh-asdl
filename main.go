package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/the-kaustubh/asdl/repository"
	"github.com/tidwall/gjson"
)

func main() {
	var userName, repoName, token, assetQuery, output string
	flag.StringVar(&userName, "user", "", "Github username")
	flag.StringVar(&repoName, "repo", "", "Github repository name")
	flag.StringVar(&token, "token", "", "Secret token key for private repos")
	flag.StringVar(&assetQuery, "asset", "", "Asset query string")
	flag.StringVar(&output, "o", "", "Output file name")
	flag.Parse()

	fmt.Println(userName, repoName, token, assetQuery, output)
	if len(token) == 0 {
		fmt.Println("Please provide valid token")
		return
	}

	repo := repository.Repo{
		Username: userName,
		Name:     repoName,
	}

	client := http.Client{}
	assets, err := getAssets(client, token, &repo)
	if err != nil {
		log.Fatal(err)
	}
	parsed := gjson.Parse(assets)
	// def := "0.assets.0.id"
	asset_id := parsed.Get(assetQuery).String()
	fmt.Println(asset_id)
	err = downloadAsset(client, &repo, token, asset_id, output)
	if err != nil {
		log.Fatal(err)
	}
}

func getAssets(client http.Client, token string, repo *repository.Repo) (string, error) {
	req, err := http.NewRequest("GET", repo.GetFullRepoURL(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("Accept", "application/vnd.github.v3.raw")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func downloadAsset(client http.Client, repo *repository.Repo, token, asset_id, filename string) error {
	url := repo.GetAssetUrlWithToken(asset_id, token)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/octet-stream")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println(size, "bytes downloaded")

	return nil
}
