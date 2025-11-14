package nodeping

import (
	"encoding/json"
	"fmt"
	"sort"

	"gopkg.in/resty.v1"
)

const (
	BaseURL = "https://api.nodeping.com/api/1"
	Version = "0.0.1"
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
	Error       NodePingError
	R           *resty.Request
	MockResults string
}

// New creates a new Client
func New(config ClientConfig) (*Client, error) {
	var client Client

	if config.Token == "" {
		return &Client{}, fmt.Errorf("token is required in ClientConfig")
	}
	client.Config.Token = config.Token

	client.Config.BaseURL = BaseURL
	if config.BaseURL != "" {
		client.Config.BaseURL = config.BaseURL
	}

	client.Config.CustomerID = config.CustomerID
	client.MockResults = ""

	resty.SetHostURL(client.Config.BaseURL)
	resty.SetBasicAuth(client.Config.Token, "")
	resty.SetHeader("user-agent", "sil-org/app-monitoring-archiver "+Version)
	client.R = resty.R()
	client.R.SetError(&client.Error)

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
		_, err := c.R.SetResult(&listObj).Get(path)
		errChk := c.CheckForError(err)
		if errChk != nil {
			return []CheckResponse{}, errChk
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

	_, err := c.R.SetResult(&check).Get(path)
	errChk := c.CheckForError(err)
	if errChk != nil {
		return CheckResponse{}, errChk
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

	_, err := c.R.SetResult(&listObj).Get(path)
	errChk := c.CheckForError(err)
	if errChk != nil {
		return map[string]UptimeResponse{}, errChk
	}

	return listObj, nil
}

// ListContactGroups retrieves the list of Contact Groups
func (c *Client) ListContactGroups() (map[string]ContactGroupResponse, error) {
	path := "/contactgroups"
	if c.Config.CustomerID != "" {
		path = fmt.Sprintf("/checks/%s", c.Config.CustomerID)
	}
	var listObj map[string]ContactGroupResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &listObj)
		return listObj, nil
	}

	_, err := c.R.SetResult(&listObj).Get(path)
	errChk := c.CheckForError(err)
	if errChk != nil {
		return map[string]ContactGroupResponse{}, errChk
	}

	return listObj, nil
}

func (c *Client) CheckForError(err error) error {
	if err != nil {
		return err
	}
	if c.Error.Error != "" {
		return fmt.Errorf("%s", c.Error.Error)
	}
	return nil
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
