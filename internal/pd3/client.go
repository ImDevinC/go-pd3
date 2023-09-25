package pd3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ImDevinC/go-pd3/internal/models"
)

const nebulaBaseUrl string = "https://nebula.starbreeze.com"
const challengesEndpoint string = "/challenge/v1/public/namespaces/pd3/users/me/records"

type Client struct {
	httpClient *http.Client
	token      string
	baseUrl    string
}

type PD3Option func(*Client)

func New(opts ...PD3Option) *Client {
	client := Client{
		httpClient: http.DefaultClient,
		baseUrl:    nebulaBaseUrl,
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

func WithHttpClient(client *http.Client) PD3Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithToken(token string) PD3Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithBaseUrl(baseUrl string) PD3Option {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

func (c *Client) GetChallenges() ([]models.PD3DataResponse, error) {
	challenges := []models.PD3DataResponse{}
	params := url.Values{}
	params.Add("limit", "100")
	reqUrl, err := url.Parse(c.baseUrl + challengesEndpoint)
	if err != nil {
		return challenges, err
	}
	reqUrl.RawQuery = params.Encode()
	nextUrl := reqUrl.String()
	for {
		req, err := http.NewRequest(http.MethodGet, nextUrl, nil)
		if err != nil {
			return challenges, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
		res, err := c.httpClient.Do(req)
		if err != nil {
			return challenges, err
		}
		if res.StatusCode != http.StatusOK {
			return challenges, errors.New(res.Status)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return challenges, err
		}
		var pd3Response models.PD3Response
		err = json.Unmarshal(body, &pd3Response)
		if err != nil {
			return challenges, err
		}
		challenges = append(challenges, pd3Response.Data...)
		if pd3Response.Paging.Next == "" {
			break
		}
		nextUrl = c.baseUrl + pd3Response.Paging.Next
	}
	return challenges, nil
}
