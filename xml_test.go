package pastebin

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestPasteXmlDecode(t *testing.T) {
	// copied from https://pastebin.com/api#9
	data := []byte(`<paste> <paste_key>0b42rwhf</paste_key> <paste_date>1297953260</paste_date> <paste_title>javascript test</paste_title> <paste_size>15</paste_size> <paste_expire_date>1297956860</paste_expire_date> <paste_private>0</paste_private> <paste_format_long>JavaScript</paste_format_long> <paste_format_short>javascript</paste_format_short> <paste_url>https://pastebin.com/0b42rwhf</paste_url> <paste_hits>15</paste_hits> </paste>`)

	var x PasteInfo
	if err := xml.Unmarshal(data, &x); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := PasteInfo{
		XMLName:    xml.Name{Local: "paste"},
		Key:        "0b42rwhf",
		CreateTS:   1297953260,
		Title:      "javascript test",
		Size:       15,
		ExpireTS:   1297956860,
		AccessMode: Public,
		Format:     "JavaScript",
		FormatCode: "javascript",
		URL:        "https://pastebin.com/0b42rwhf",
		Hits:       15,
	}

	if !reflect.DeepEqual(expected, x) {
		t.Fatalf("unexpected result: %+v", x)
	}
}

func TestUserXmlDecode(t *testing.T) {
	// copied from https://pastebin.com/api#12
	data := []byte(`<user> <user_name>wiz_kitty</user_name> <user_format_short>text</user_format_short> <user_expiration>N</user_expiration> <user_avatar_url>https://pastebin.com/cache/a/1.jpg</user_avatar_url> <user_private>1</user_private><user_website>https://myawesomesite.com</user_website> <user_email>oh@dear.com</user_email> <user_location>New York</user_location> <user_account_type>1</user_account_type> </user>`)

	var x UserInfo
	if err := xml.Unmarshal(data, &x); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := UserInfo{
		XMLName:     xml.Name{Local: "user"},
		Name:        "wiz_kitty",
		FormatCode:  "text",
		Expiration:  Never,
		AvatarURL:   "https://pastebin.com/cache/a/1.jpg",
		AccessMode:  Unlisted,
		Website:     "https://myawesomesite.com",
		EMail:       "oh@dear.com",
		Location:    "New York",
		AccountType: 1,
	}

	if !reflect.DeepEqual(expected, x) {
		t.Fatalf("unexpected result: %+v", x)
	}
}
