package pastebin

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"strconv"
	"time"
)

// AccessMode indicates who can access this paste
type AccessMode int

func (m AccessMode) String() string {
	return strconv.Itoa(int(m))
}

// possible modes
const (
	Public AccessMode = iota
	Unlisted
	Private
)

// PasteInfo represents the info of a paste retrieved from pastebin
type PasteInfo struct {
	XMLName    xml.Name   `xml:"paste"`
	Key        string     `xml:"paste_key"`
	CreateTS   int64      `xml:"paste_date"`
	Title      string     `xml:"paste_title"`
	Size       int64      `xml:"paste_size"`
	ExpireTS   int64      `xml:"paste_expire_date"`
	AccessMode AccessMode `xml:"paste_private"`
	Format     string     `xml:"paste_format_long"`  // detailed name of hightlight format
	FormatCode string     `xml:"paste_format_short"` // code of hightlight format
	URL        string     `xml:"paste_url"`
	Hits       int64      `xml:"paste_hits"`
}

func (p *PasteInfo) fromTS(ts int64) time.Time {
	return time.Unix(
		ts/1000,
		(ts%1000)*int64(time.Millisecond),
	)
}

// CreateAt converts CreateTS to time.Time
func (p *PasteInfo) CreateAt() time.Time {
	return p.fromTS(p.CreateTS)
}

// ExpireAt converts ExpireTS to time.Time
func (p *PasteInfo) ExpireAt() time.Time {
	return p.fromTS(p.ExpireTS)
}

// Expiration represents parameter of specifying expiration date
type Expiration string

// supported formats
const (
	Never Expiration = "N"
	In10M            = "10M"
	In1H             = "1H"
	In1D             = "1D"
	In1W             = "1W"
	In2W             = "2W"
	In1M             = "1M"
	In6M             = "6M"
	In1Y             = "1Y"
)

// Paste represents a new paste
type Paste struct {
	Title      string
	Content    string
	AccessMode AccessMode
	Format     string
	ExpireAt   Expiration
	UserKey    string
}

func (p *Paste) Values() (ret url.Values) {
	ret = url.Values{}
	ret.Set("api_paste_code", p.Content)
	ret.Set("api_paste_private", p.AccessMode.String())
	if p.Title != "" {
		ret.Set("api_paste_name", p.Title)
	}
	if x := string(p.ExpireAt); x != "" {
		ret.Set("api_paste_expire_date", x)
	}

	if p.Format != "" {
		ret.Set("api_paste_format", p.Format)
	}
	if p.UserKey != "" {
		ret.Set("api_user_key", p.UserKey)
	}

	return
}

// Error indicates the remote endpoint returns some error
type Error string

func (e Error) Error() string {
	return string(e)
}

func isError(data []byte) (ret string, err error) {
	if bytes.HasPrefix(data, []byte(`Bad API request, `)) {
		err = Error(string(data))
		return
	}
	ret = string(data)
	return
}

// UserInfo represents the info of a user retrieved from pastebin
type UserInfo struct {
	XMLName     xml.Name   `xml:"user"`
	Name        string     `xml:"user_name"`
	FormatCode  string     `xml:"user_format_short"`
	Expiration  Expiration `xml:"user_expiration"`
	AvatarURL   string     `xml:"user_avatar_url"`
	AccessMode  AccessMode `xml:"user_private"`
	Website     string     `xml:"user_website"`
	EMail       string     `xml:"user_email"`
	Location    string     `xml:"user_location"`
	AccountType int        `xml:"user_account_type"` // 0 normal, 1 PRO
}
