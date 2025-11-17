package nodeping

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

const (
	DefaultBaseURL = "https://api.nodeping.com/api/1"
	Version        = "0.0.1"
)

// ClientConfig type includes configuration options for NodePing client.
type ClientConfig struct {
	BaseURL    string
	Token      string
	CustomerID string
}

// Client holds config and provides methods for various api calls
type Client struct {
	Config      ClientConfig
	httpClient  *http.Client
	MockResults string
}

// New creates a new Client
func New(config ClientConfig) (*Client, error) {
	client := Client{Config: config}

	if config.Token == "" {
		return nil, fmt.Errorf("token is required in ClientConfig")
	}

	if config.BaseURL == "" {
		client.Config.BaseURL = DefaultBaseURL
	}

	client.httpClient = &http.Client{Timeout: time.Second * 30}

	return &client, nil
}

// ListChecks retrieves all the "Checks" in NodePing
func (c *Client) ListChecks() ([]CheckResponse, error) {
	path := "/checks"
	if c.Config.CustomerID != "" {
		path = fmt.Sprintf("/checks/%s", c.Config.CustomerID)
	}
	var listObj map[string]CheckResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &listObj)
	} else {
		if err := c.request(path, &listObj); err != nil {
			return nil, err
		}
	}

	var list []CheckResponse
	for _, item := range listObj {
		list = append(list, item)
	}

	return list, nil
}

// GetCheck retrieves data about one Check using its id
func (c *Client) GetCheck(id string) (CheckResponse, error) {
	path := fmt.Sprintf("/checks/%s", id)
	var check CheckResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &check)
		return check, nil
	}

	if err := c.request(path, &check); err != nil {
		return CheckResponse{}, err
	}

	return check, nil
}

// GetUptime retrieves the uptime entries for a certain check within an optional date range (by Timestamp with microseconds)
func (c *Client) GetUptime(id string, period Period) (map[string]UptimeResponse, error) {
	// Build path with potentially a "?" and "&" symbols
	// I'm not getting c.R's built in query param functions to work properly
	pathDelimiter := ""
	queryParams := ""
	queryParamDelimiter := ""

	if !period.From.IsZero() {
		queryParams = fmt.Sprintf("start=%d", period.From.Unix()*1000)
		queryParamDelimiter = "&"
		pathDelimiter = "?"
	}

	if !period.To.IsZero() {
		queryParams = fmt.Sprintf("%s%send=%d", queryParams, queryParamDelimiter, period.To.Unix()*1000)
		pathDelimiter = "?"
	}

	path := fmt.Sprintf("/results/uptime/%s%s%s", id, pathDelimiter, queryParams)

	var listObj map[string]UptimeResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &listObj)
		return listObj, nil
	}

	if err := c.request(path, &listObj); err != nil {
		return nil, err
	}

	return listObj, nil
}

// ListContactGroups retrieves the list of Contact Groups
func (c *Client) ListContactGroups() (map[string]ContactGroupResponse, error) {
	path := "/contactgroups"
	if c.Config.CustomerID != "" {
		path = fmt.Sprintf("/contactgroups/%s", c.Config.CustomerID)
	}
	var listObj map[string]ContactGroupResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &listObj)
		return listObj, nil
	}
	if err := c.request(path, &listObj); err != nil {
		return nil, err
	}

	return listObj, nil
}

func (c *Client) GetContactGroupIDFromName(contactGroupName string) (string, error) {
	contactGroups, err := c.ListContactGroups()
	if err != nil {
		return "", fmt.Errorf("error retrieving contact groups: %w", err)
	}

	cgID := ""
	for cgKey, cg := range contactGroups {
		if cg.Name == contactGroupName {
			cgID = cgKey
			break
		}
	}

	if cgID == "" {
		return "", fmt.Errorf(`contact group not found with name: "%s"`, contactGroupName)
	}

	return cgID, nil
}

func (c *Client) GetCheckIDsAndLabels(id string) ([]string, map[string]string, error) {
	checkIDs := map[string]string{}
	checkLabels := []string{}

	checks, err := c.ListChecks()
	if err != nil {
		return checkLabels, checkIDs, err
	}

	for _, check := range checks {
		// Notifications is a list of maps with the contactGroup ID as keys
		for _, notification := range check.Notifications {
			foundOne := false
			for nKey := range notification {
				if nKey == id {
					checkIDs[check.Label] = check.ID
					checkLabels = append(checkLabels, check.Label)
					foundOne = true
					break
				}
			}
			if foundOne {
				break
			}
		}
	}

	sort.Strings(checkLabels)
	return checkLabels, checkIDs, nil
}

func (c *Client) GetUptimesForChecks(checkIDs map[string]string, period Period) map[string]float32 {
	uptimes := map[string]float32{}

	for _, checkID := range checkIDs {
		nextUptime, err := c.GetUptime(checkID, period)
		if err != nil {
			fmt.Printf("Error getting uptime for check ID %s.\n%s\n", checkID, err.Error())
			continue
		}
		uptimes[checkID] = nextUptime["total"].Uptime
	}

	return uptimes
}

func (c *Client) request(path string, v any) error {
	req, err := http.NewRequest("GET", c.Config.BaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("user-agent", "sil-org/app-monitoring-archiver "+Version)

	req.SetBasicAuth(c.Config.Token, "")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d, body: %s", res.StatusCode, body[0:min(250, len(body))])
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("invalid response body %s: %w", body, err)
	}
	return nil
}

func GetUptimesForContactGroup(token, group string, period Period) (UptimeResults, error) {
	var emptyResults UptimeResults
	npClient, err := New(ClientConfig{Token: token})
	if err != nil {
		return emptyResults, fmt.Errorf("error initializing cli: %w", err)
	}

	cgID, err := npClient.GetContactGroupIDFromName(group)
	if err != nil {
		return emptyResults, err
	}

	checkLabels, checkIDs, err := npClient.GetCheckIDsAndLabels(cgID)
	if err != nil {
		return emptyResults, err
	}

	uptimes := npClient.GetUptimesForChecks(checkIDs, period)
	uptimesByLabel := map[string]float32{}

	for _, label := range checkLabels {
		uptimesByLabel[label] = uptimes[checkIDs[label]]
	}

	results := UptimeResults{
		CheckLabels: checkLabels,
		Uptimes:     uptimesByLabel,
		StartTime:   period.From.Unix(),
		EndTime:     period.To.Unix(),
	}

	return results, nil
}
