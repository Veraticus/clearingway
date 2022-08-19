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
	"strings"
	"time"

	"github.com/Veraticus/clearingway/internal/ffxiv"

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

func (f *Fflogs) GetEncounterRankings(encounterIds []int, char *ffxiv.Character) (*EncounterRankings, error) {
	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("query{characterData{character(name: \"%s\", serverSlug: \"%s\", serverRegion: \"NA\"){", char.Name, char.Server))
	for _, encounterId := range encounterIds {
		query.WriteString(fmt.Sprintf("e%d: encounterRankings(encounterID: %d)", encounterId, encounterId))
	}
	query.WriteString("}}}")

	raw, err := f.graphqlClient.ExecRaw(context.Background(), query.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %w", err)
	}

	var characterData map[string]*json.RawMessage
	err = json.Unmarshal(raw, &characterData)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	var character map[string]*json.RawMessage
	err = json.Unmarshal(*characterData["characterData"], &character)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}
	if character["character"] == nil {
		return nil, fmt.Errorf("Character %s (%s) not found in fflogs!", char.Name, char.Server)
	}

	var rawEncounters map[string]*json.RawMessage
	err = json.Unmarshal(*character["character"], &rawEncounters)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	encounterRankings := &EncounterRankings{Encounters: map[int]*EncounterRanking{}}
	for rawId, rawEncounter := range rawEncounters {
		var id int
		_, err = fmt.Sscanf(rawId, "e%d", &id)
		if err != nil {
			return nil, fmt.Errorf("Could not get encounter ID back from response: %w", err)
		}

		encounterRanking := &EncounterRanking{}
		err = json.Unmarshal(*rawEncounter, encounterRanking)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
		}

		encounterRankings.Encounters[id] = encounterRanking
	}

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
