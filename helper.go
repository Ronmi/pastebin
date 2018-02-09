package pastebin

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

var reDevKey *regexp.Regexp

const (
	devKeyPrefix = `<div class="code_box">`
	devKeySuffix = `</div>`
)

func init() {
	reDevKey = regexp.MustCompile(
		devKeyPrefix + `[0-9a-f]{32}` + devKeySuffix,
	)
}

// GetDevKey retrieves dev key from pastebin api page
//
// It creates a new http.Client with correct CookieJar if client == nil
func GetDevKey(acct, pass string, client *http.Client) (key string, err error) {
	if client == nil {
		client = &http.Client{}
		jar, err := cookiejar.New(nil)
		if err != nil {
			return "", err
		}

		client.Jar = jar
	}

	vals := url.Values{}
	vals.Set("user_name", acct)
	vals.Set("user_password", pass)
	vals.Set("submit_hidden", "submit_hidden")

	resp, err := client.PostForm("https://pastebin.com/login.php", vals)
	if err != nil {
		return
	}

	if resp.StatusCode >= 400 {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		return key, errors.New(resp.Status)
	}

	cookies := resp.Cookies()
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	req, err := http.NewRequest("GET", "https://pastebin.com/api", nil)
	if err != nil {
		return
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	keyHTML := reDevKey.Find(data)
	if keyHTML == nil {
		return key, errors.New("no dev key found")
	}

	key = string(keyHTML[len(devKeyPrefix):])
	key = key[:len(key)-len(devKeySuffix)]
	return
}
