package kstartup

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const announcementPath = "https://www.k-startup.go.kr/common/announcement/announcementList.do"

type CannotNotFoundElement struct {
	selector string
}

func (err *CannotNotFoundElement) Error() string {
	return err.selector
}

type Client struct {
	http *http.Client
}

func New() *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	var httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: jar,
	}
	return &Client{http: httpClient}
}

func createRequestBody(csrf string) url.Values {
	return url.Values{
		"CSRF_NONCE":      []string{csrf},
		"mid":             []string{"30004"},
		"bid":             []string{"701"},
		"pageIndex":       []string{"37"},
		"searchSortOrder": []string{"INSERT_DATE"},
	}
}

func (client *Client) getCsrfToken() (string, error) {
	res, err := client.http.Get(announcementPath)
	if err != nil {
		return "", err
	}

	defer func() {
		err = res.Body.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	selector := "input[name=CSRF_NONCE]"
	selection := doc.Find(selector)

	csrf, exist := selection.Attr("value")
	if !exist {
		return "", &CannotNotFoundElement{selector: selector}
	}

	return csrf, nil

}

func (client *Client) post(url string, value url.Values) (*http.Response, error) {
	return client.http.Post(
		url,
		"application/x-www-form-urlencoded",
		strings.NewReader(value.Encode()),
	)
}

func trim(text string) string {
	return strings.TrimSpace(text)
}

func (client *Client) GetAnnouncements() ([]Announcement, error) {
	csrf, err := client.getCsrfToken()
	if err != nil {
		return nil, err
	}

	value := createRequestBody(csrf)
	res, err := client.post(announcementPath, value)

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	annList := []Announcement{}

	annListSelection := doc.Find(".ann_list")
	annListSelection.Children().Each(func(i int, s *goquery.Selection) {
		title, exist := s.Find("a").Attr("title")
		if !exist {
			title = ""
		}
		annList = append(annList, Announcement{
			ID:           1,
			BoardID:      "",
			Group:        "",
			Title:        trim(title),
			Name:         "",
			Organization: "",
			Due:          time.Now(),
		})
	})
	log.Println(annList)
	return nil, nil
}
