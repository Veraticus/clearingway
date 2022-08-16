package fflogs

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
)

type Fflogs struct {
	clientId      string
	clientSecret  string
	graphqlClient *graphql.Client
}

type fflogsAccessToken struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	AccessToken      string `json:"access_token"`
}

var encounterRankingsQuery struct {
	CharacterData struct {
		Character struct {
			EncounterRankings json.RawMessage `graphql:"encounterRankings(encounterID: $encounterId)"`
		} `graphql:"character(name: $characterName, serverSlug: $characterServer, serverRegion: \"NA\")"`
	} `graphql:"characterData"`
}

func (f *Fflogs) GetEncounterRankings(encounterId int, characterName, characterServer string) (*EncounterRankings, error) {
	variables := map[string]interface{}{
		"encounterId":     graphql.Int(encounterId),
		"characterName":   graphql.String(characterName),
		"characterServer": graphql.String(characterServer),
	}

	err := f.graphqlClient.Query(context.Background(), &encounterRankingsQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("Could not execute graphql query: %w", err)
	}

	encounterRankings := &EncounterRankings{}
	err = json.Unmarshal(encounterRankingsQuery.CharacterData.Character.EncounterRankings, encounterRankings)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	fmt.Printf("Results are: %+v\n", encounterRankings)

	return encounterRankings, nil
}

func (f *Fflogs) SetGraphqlClient() {
	src := oauth2.ReuseTokenSource(nil, f)
	httpClient := oauth2.NewClient(context.Background(), src)
	f.graphqlClient = graphql.NewClient("https://www.fflogs.com/api/v2/client", httpClient)
}

func (f *Fflogs) Token() (*oauth2.Token, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	grantTypeField, err := writer.CreateFormField("grant_type")
	if err != nil {
		return nil, fmt.Errorf("Could not create form field: %w", err)
	}
	_, err = grantTypeField.Write([]byte("client_credentials"))
	if err != nil {
		return nil, fmt.Errorf("Could not write form field: %w", err)
	}

	req, err := http.NewRequest("POST", "https://www.fflogs.com/oauth/token", body)
	if err != nil {
		return nil, fmt.Errorf("Could not create new HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(f.clientId, f.clientSecret))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP oauth2 token request failed: %w", err)
	}
	defer resp.Body.Close()

	returnedFflogsToken := &fflogsAccessToken{}
	err = json.NewDecoder(resp.Body).Decode(returnedFflogsToken)
	if err != nil {
		return nil, fmt.Errorf("Could not coerce response to JSON: %w", err)
	}

	if returnedFflogsToken.Error != "" {
		return nil, fmt.Errorf("Token error %v: %v", returnedFflogsToken.Error, returnedFflogsToken.ErrorDescription)
	}

	token := &oauth2.Token{}
	token.AccessToken = returnedFflogsToken.AccessToken
	token.TokenType = returnedFflogsToken.TokenType
	token.Expiry = time.Now().Local().Add(time.Second * time.Duration(returnedFflogsToken.ExpiresIn))

	return token, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
