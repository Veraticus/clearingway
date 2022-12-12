package fflogs

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

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

type RankingToGet struct {
	IDs        []int
	Difficulty int
}

func (f *Fflogs) SetCharacterLodestoneID(char *ffxiv.Character) error {
	query := fmt.Sprintf(
		"query{characterData{character(name: \"%s\", serverSlug: \"%s\", serverRegion: \"%s\"){lodestoneID}}}",
		char.Name(),
		char.World,
		char.PhysicalDatacenter().Abbreviation,
	)

	raw, err := f.graphqlClient.ExecRaw(context.Background(), query, nil)
	if err != nil {
		return fmt.Errorf("Error executing query: %w", err)
	}

	var characterData map[string]*json.RawMessage
	err = json.Unmarshal(raw, &characterData)
	if err != nil {
		return fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	var character map[string]*json.RawMessage
	err = json.Unmarshal(*characterData["characterData"], &character)
	if err != nil {
		return fmt.Errorf("Could not unmarshal JSON: %w", err)
	}
	if character["character"] == nil {
		return fmt.Errorf("Character %s (%s) not found in fflogs!", char.Name(), char.World)
	}

	var rawCharacterResponse map[string]*json.RawMessage
	err = json.Unmarshal(*character["character"], &rawCharacterResponse)
	if err != nil {
		return fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	for rawKey, rawValue := range rawCharacterResponse {
		if rawKey == "lodestoneID" {
			var id int
			err = json.Unmarshal(*rawValue, &id)
			if err != nil {
				return fmt.Errorf("Could not unmarshal lodestone ID: %w", err)
			}

			char.LodestoneID = id
			return nil
		}
	}

	return fmt.Errorf("Lodestone ID not found on fflogs!")
}

var returnedRankingsRegexp = regexp.MustCompile(`(\D+)Z(\d+)`)

func (f *Fflogs) GetRankingsForCharacter(rankingsToGet []*RankingToGet, char *ffxiv.Character) (*Rankings, error) {
	query := strings.Builder{}
	query.WriteString(
		fmt.Sprintf(
			"query{characterData{character(name: \"%s\", serverSlug: \"%s\", serverRegion: \"%s\"){",
			char.Name(),
			char.World,
			char.PhysicalDatacenter().Abbreviation,
		),
	)
	for _, rankingToGet := range rankingsToGet {
		for _, id := range rankingToGet.IDs {
			query.WriteString(
				fmt.Sprintf(
					"rdpsZ%d:encounterRankings(encounterID: %d, difficulty: %d, metric: rdps) ",
					id,
					id,
					rankingToGet.Difficulty,
				),
			)
			query.WriteString(
				fmt.Sprintf(
					"hpsZ%d:encounterRankings(encounterID: %d, difficulty: %d, metric: hps) ",
					id,
					id,
					rankingToGet.Difficulty,
				),
			)
		}
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
		return nil, fmt.Errorf("Character %s (%s) not found in fflogs!", char.Name(), char.World)
	}

	var rawRankings map[string]*json.RawMessage
	err = json.Unmarshal(*character["character"], &rawRankings)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
	}

	rankings := &Rankings{Rankings: map[int]*Ranking{}}
	for rawId, rawRanking := range rawRankings {
		match := returnedRankingsRegexp.FindStringSubmatch(rawId)
		if match == nil {
			return nil, fmt.Errorf("Returned stanza did not match expected format: %v\n", rawId)
		}

		metric := match[1]
		idString := match[2]
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, fmt.Errorf("Could not convert id %v from string to int: %v\n", idString, err)
		}

		ranking := &Ranking{Metric: Metric(metric)}
		err = json.Unmarshal(*rawRanking, ranking)
		if err != nil {
			return nil, fmt.Errorf("Could not unmarshal JSON: %w", err)
		}
		if ranking.Error != "" {
			if ranking.Error == "Invalid encounter id specified." {
				fmt.Printf("Could not find encounters for id %d, continuing...\n", id)
				continue
			} else {
				return nil, fmt.Errorf("Received error from fflogs for encounter %d: %v", id, ranking.Error)
			}
		}

		rankings.Add(id, ranking)
	}

	return rankings, nil
}

func (f *Fflogs) SetGraphqlClient() {
	src := oauth2.ReuseTokenSource(nil, f)
	httpClient := oauth2.NewClient(context.Background(), src)
	newClient := graphql.NewClient("https://www.fflogs.com/api/v2/client", httpClient)
	newClient = newClient.WithRequestModifier(func(req *http.Request) {
		req.Header.Add("User-Agent", "clearingway")
	})
	f.graphqlClient = newClient
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
