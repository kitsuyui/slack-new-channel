package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
)

// SlackChannelAPIResponse represents response of this API:
// https://api.slack.com/methods/channels.list
type SlackChannelAPIResponse struct {
	Ok           bool     `json:"ok"`
	Channels     Channels `json:"channels"`
	ErrorMessage string   `json:"error"`
}

// Channel represents a channel objects in SlackChannelAPIResponse.
// https://api.slack.com/types/channel
type Channel struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	IsChannel     bool        `json:"is_channel"`
	Created       json.Number `json:"created,Number"`
	Creator       string      `json:"creator"`
	IsArchived    bool        `json:"is_archived"`
	IsGeneral     bool        `json:"is_general"`
	IsMember      bool        `json:"is_member"`
	Members       []string    `json:"members"`
	Topic         Topic       `json:"topic"`
	Purpose       Purpose     `json:"purpose"`
	PreviousNames []string    `json:"previous_names"`
	NumMembers    json.Number `json:"num_members,Number"`
}

// Topic represents topic item of channel object.
type Topic struct {
	Value   string      `json:"value"`
	Creator string      `json:"creator"`
	LastSet json.Number `json:"last_set,Number"`
}

// Purpose represents purpose item of channel object.
type Purpose struct {
	Value   string      `json:"value"`
	Creator string      `json:"creator"`
	LastSet json.Number `json:"last_set,Number"`
}

// Channels are a list of channel.
type Channels []Channel

// ByCreated wraps Channels for it can be sorted by its created date.
type ByCreated struct {
	Channels
}

// Swap is required for sorting.
func (p Channels) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Len is required for sorting.
func (b ByCreated) Len() int {
	return len(b.Channels)
}

// Less is required for sorting.
func (b ByCreated) Less(i, j int) bool {
	return b.Channels[i].Created < b.Channels[j].Created
}

// GetSlackChannelAPIResponse simply fetch Slack channels.list API and returns it result.
// Slack access token are required.
func GetSlackChannelAPIResponse(token string) (d *SlackChannelAPIResponse, err error) {
	requrl := "https://slack.com/api/channels.list"
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("token", token)
	query.Add("exclude_archived", "true")
	req.URL.RawQuery = query.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// LoadChannelFromFile loads channel object in file.
func LoadChannelFromFile(path string) (c *Channel, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		log.Fatal(err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// SlackPost just comments in Slack with Incoming Webhooks
// https://api.slack.com/incoming-webhooks
func SlackPost(username string, icon string, text string, hookurl string) (err error) {
	apiURL, err := url.ParseRequestURI(hookurl)
	if err != nil {
		return err
	}
	query := url.Values{}
	apiURL.RawQuery = query.Encode()
	data, _ := json.Marshal(map[string]string{
		"text":       text,
		"icon_emoji": icon,
		"username":   username,
	})
	client := &http.Client{}
	r, err := http.NewRequest("POST", hookurl, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	_, err = client.Do(r)
	if err != nil {
		return err
	}
	return nil
}

func main2(token string, path string, hookurl string) {
	d, err := GetSlackChannelAPIResponse(token)
	if err != nil {
		log.Fatal(err)
	}
	if !d.Ok {
		log.Fatal(d.ErrorMessage)
	}
	created := json.Number(0)
	maxCh, err := LoadChannelFromFile(path)

	if err == nil {
		created = maxCh.Created
	}

	channels := d.Channels
	sort.Sort(ByCreated{channels})

	s := ""
	for _, ch := range d.Channels {
		if created < ch.Created && !ch.IsArchived {
			s += "<#" + ch.ID + "|" + ch.Name + ">\n"
		}
	}

	if s != "" {
		SlackPost("New Channel Report", ":new:", s, hookurl)
	}

	newMaxCh := channels[len(channels)-1]
	b, err := json.Marshal(newMaxCh)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(path, b, 0644)
}

var interval = flag.Duration("interval", 600*time.Second, "Interval duration.")

func init() {
	flag.DurationVar(interval, "i", 600*time.Second, "Interval duration.")
}

func main() {
	token := os.Getenv("SLACK_API_TOKEN")
	path := os.Getenv("LATEST_CHANNEL_JSON_PATH")
	hookurl := os.Getenv("SLACK_WEBHOOK_URL")

	usage := `slack-new-channel
Usage:
 slack-new-channel [--daemon]
Options:
 -d --daemon  Daemon mode.
 -i=<duration> --interval=<duration>  Interval [default: 600s].
 -h --help  Show this screen.
`

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Errorf("Err: %s", err)
	}

	timeoutStr, err := opts.String("--interval")
	if err != nil {
		fmt.Errorf("%s", err)
	}

	var duration time.Duration
	f, err := strconv.ParseFloat(timeoutStr, 64)
	if err != nil {
		duration, err = time.ParseDuration(timeoutStr)
		if err != nil {
			fmt.Errorf("%s", err)
		}
	} else {
		duration = time.Duration(f) * time.Second
	}

	daemonMode, _ := opts.Bool("--daemon")
	for daemonMode {
		main2(token, path, hookurl)
		time.Sleep(duration)
	}
}
