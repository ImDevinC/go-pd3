package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ImDevinC/go-pd3/models"
	_ "github.com/joho/godotenv/autoload"
)

const baseUrl string = "https://nebula.starbreeze.com"
const firstUrl string = "/challenge/v1/public/namespaces/pd3/users/me/records"

var token = os.Getenv("NEBULA_BEARER_TOKEN")

var outputFile string
var flagToken string

func init() {
	flag.StringVar(&outputFile, "outputFile", "", "Filename for results to be written to")
	flag.StringVar(&flagToken, "token", "", "Bearer token to use")
}

func main() {
	flag.Parse()
	if token == "" && flagToken == "" {
		log.Fatal("Missing token. Set NEBULA_BEARER_TOKEN or use -token flag")
	}
	if flagToken != "" {
		token = flagToken
	}
	params := url.Values{}
	params.Add("limit", "100")
	reqUrl, err := url.Parse(baseUrl + firstUrl)
	if err != nil {
		log.Fatal(err)
	}
	reqUrl.RawQuery = params.Encode()
	challenges := []models.PD3DataResponse{}
	nextUrl := reqUrl.String()
	for {
		resp, next, err := getChallenges(nextUrl, token)
		if err != nil {
			log.Fatal(err)
		}
		challenges = append(challenges, resp...)
		if next != "" {
			nextUrl = baseUrl + next
			continue
		}
		break
	}
	if outputFile == "" {
		log.Printf("%+v\n", challenges)
		return
	}
	data, err := json.Marshal(challenges)
	if err != nil {
		log.Fatal(err)
	}
	prettyData := bytes.Buffer{}
	err = json.Indent(&prettyData, data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(outputFile, prettyData.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getChallenges(finalUrl string, token string) ([]models.PD3DataResponse, string, error) {
	req, err := http.NewRequest(http.MethodGet, finalUrl, nil)
	if err != nil {
		return []models.PD3DataResponse{}, "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []models.PD3DataResponse{}, "", err
	}
	if res.StatusCode != http.StatusOK {
		return []models.PD3DataResponse{}, "", errors.New(res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []models.PD3DataResponse{}, "", err
	}
	var pd3Response models.PD3Response
	err = json.Unmarshal(body, &pd3Response)
	if err != nil {
		return []models.PD3DataResponse{}, "", err
	}
	return pd3Response.Data, pd3Response.Paging.Next, nil
}
