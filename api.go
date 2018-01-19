package pastebin

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// endpoints
const (
	_post  = `https://pastebin.com/api/api_post.php`
	_login = `https://pastebin.com/api/api_login.php`
	_raw   = `https://pastebin.com/api/api_raw.php`
)

// API warps pastebin api for you
//
// It is safe to share between multiple goroutines.
type API struct {
	Client *http.Client // use http.DefaultClient if nil
	Key    string       // api_dev_key
}

func (a *API) client() (c *http.Client) {
	c = a.Client
	if c == nil {
		c = http.DefaultClient
	}

	return
}

func (a *API) post(val url.Values) (resp *http.Response, err error) {
	val.Set("api_dev_key", a.Key)
	return a.client().PostForm(_post, val)
}

func (a *API) raw(val url.Values) (resp *http.Response, err error) {
	val.Set("api_dev_key", a.Key)
	return a.client().PostForm(_raw, val)
}

func (a *API) login(val url.Values) (resp *http.Response, err error) {
	val.Set("api_dev_key", a.Key)
	return a.client().PostForm(_login, val)
}

// UserKey allocates an api_user_key using specified user info
func (a *API) UserKey(acct, pass string) (key string, err error) {
	val := url.Values{}
	val.Set("api_user_name", acct)
	val.Set("api_user_password", pass)
	resp, err := a.login(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	key, err = isError(bytes.TrimSpace(data))
	return
}

// Post posts a new paste
func (a *API) Post(p *Paste) (uri string, err error) {
	val := p.Values()
	val.Set("api_option", "paste")
	resp, err := a.post(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	uri, err = isError(bytes.TrimSpace(data))
	return
}

func (a *API) decode(resp *http.Response) (ret []*PasteInfo, err error) {
	dec := xml.NewDecoder(resp.Body)
	ret = make([]*PasteInfo, 0)
	for err == nil {
		var info PasteInfo
		if err = dec.Decode(&info); err == nil {
			ret = append(ret, &info)
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

// List lists pastes created by the user
//
// Pass limit < 1 will be forced to 50, and 1000 if > 1000
func (a *API) List(userKey string, limit int) (ret []*PasteInfo, err error) {
	if limit < 1 {
		limit = 50
	} else if limit > 1000 {
		limit = 1000
	}

	val := url.Values{}
	val.Set("api_user_key", userKey)
	val.Set("api_results_limit", strconv.Itoa(limit))
	val.Set("api_option", "list")
	resp, err := a.post(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	return a.decode(resp)
}

// Trends lists trending pastes
func (a *API) Trends() (ret []*PasteInfo, err error) {
	val := url.Values{}
	val.Set("api_option", "trends")
	resp, err := a.post(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	return a.decode(resp)
}

// Delete deletes a paste created by user
func (a *API) Delete(userKey, pasteKey string) (err error) {
	val := url.Values{}
	val.Set("api_user_key", userKey)
	val.Set("api_paste_key", pasteKey)
	val.Set("api_option", "delete")
	resp, err := a.post(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	_, err = isError(data)
	return
}

// UserInfo retrieves info of the user
func (a *API) UserInfo(userKey string) (info *UserInfo, err error) {
	val := url.Values{}
	val.Set("api_user_key", userKey)
	val.Set("api_option", "userdetails")
	resp, err := a.post(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var i UserInfo
	if err = xml.Unmarshal(data, &i); err != nil {
		return
	}
	info = &i

	return
}

// UserPaste gets raw content of a paste created by the user
func (a *API) UserPaste(userKey, pasteKey string) (data []byte, err error) {
	val := url.Values{}
	val.Set("api_user_key", userKey)
	val.Set("api_paste_key", pasteKey)
	val.Set("api_option", "show_paste")
	resp, err := a.raw(val)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	return
}

// PubPaste gets raw content of a public/unlisted paste
func (a *API) PubPaste(pasteKey string) (data []byte, err error) {
	resp, err := a.client().Get("https://pastebin.com/raw/" + pasteKey)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	return
}
