package crawling

import (
	"encoding/xml"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type Write struct {
	Day        string
	AuthorName string
	Text       string
}

type EntryData struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"text"`
}

type ViewEntry struct {
	Key   string      `xml:"unid,attr"`
	Value []EntryData `xml:"entrydata"`
}
type ViewEntries struct {
	XMLName     xml.Name    `xml:viewentries`
	ViewEntries []ViewEntry `xml:"viewentry"`
}

func GetEvent() map[string]string {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	parsed := new(ViewEntries)
	_, body, _ := gorequest.New().Get(
		"http://ione.interpark.com/gw/app/bult/bbs00000.nsf/wviwnotice?ReadViewEntries&start=1&count=14&restricttocategory=03&page=1&&_=1504081645868",
	).Type("xml").AddCookie(
		&http.Cookie{Name: "LtpaToken", Value: "AAECAzU5QTY3Njg4NTlBN0M4MDhDTj0Ruc4RvcIRseIvT1U9MjAxMTAzMzQvTz1pbnRlcnBhcms/jQWHU+jsjHSmyRqoj3Goj/z8Qg=="},
	).End()

	_ = xml.Unmarshal([]byte(body), &parsed)

	var event Write
	var eventlist []Write

	for _, v := range parsed.ViewEntries {
		var entrydata []EntryData
		entrydata = v.Value

		for key, val := range entrydata {

			if event.AuthorName != "" && event.Day != "" && event.Text != "" {
				eventlist = append(eventlist, event)
				event.AuthorName = ""
				event.Day = ""
				event.Text = ""
			}

			switch key {
			case 1:
				event.Day = val.Value
			case 2:
				event.AuthorName = val.Value
			case 3:
				event.Text = val.Value
			}

		}
	}

	returnlist := make(map[string]string)
	var loop = 0

	for _, v := range eventlist {
		if loop < 3 {
			returnlist[v.Text] = v.Day + " " + v.AuthorName
			loop++
		}
	}

	return returnlist
}
