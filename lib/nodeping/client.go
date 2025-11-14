package nodeping

import (
	"encoding/json"
	"fmt"

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

var client Client

// Client holds config and provides methods for various api calls
type Client struct {
	Config      ClientConfig
	Error       NodePingError
	R           *resty.Request
	MockResults string
}

// Initialize new Client
func New(config ClientConfig) (*Client, error) {
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
		errChk := CheckForError(err, c)
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
	errChk := CheckForError(err, c)
	if errChk != nil {
		return CheckResponse{}, errChk
	}

	return check, nil
}

// GetUptime retrieves the uptime entries for a certain check within an optional date range (by Timestamp with microseconds)
func (c *Client) GetUptime(id string, start, end int64) (map[string]UptimeResponse, error) {
	// Build path with potentially a "?" and "&" symbols
	// I'm not getting c.R's built in query param functions to work properly
	pathDelimiter := ""
	queryParams := ""
	queryParamDelimiter := ""

	if start > 0 {
		queryParams = fmt.Sprintf("start=%d", start)
		queryParamDelimiter = "&"
		pathDelimiter = "?"
	}

	if end > 0 {
		queryParams = fmt.Sprintf("%s%send=%d", queryParams, queryParamDelimiter, end)
		pathDelimiter = "?"
	}

	path := fmt.Sprintf("/results/uptime/%s%s%s", id, pathDelimiter, queryParams)

	var listObj map[string]UptimeResponse

	if c.MockResults != "" {
		json.Unmarshal([]byte(c.MockResults), &listObj)
		return listObj, nil
	}

	_, err := c.R.SetResult(&listObj).Get(path)
	errChk := CheckForError(err, c)
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
	errChk := CheckForError(err, c)
	if errChk != nil {
		return map[string]ContactGroupResponse{}, errChk
	}

	return listObj, nil
}

func CheckForError(err error, client *Client) error {
	if err != nil {
		return err
	}
	if client.Error.Error != "" {
		return fmt.Errorf("%s", client.Error.Error)
	}
	return nil
}
