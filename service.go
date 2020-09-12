package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type series struct {
	ID         string `json:"id"`
	CnName     string `json:"cnName"`
	Poster     string `json:"poster"`
	EnName     string `json:"enName,omitempty"`
	Link       string `json:"link,omitempty"`
	RssLink    string `json:"rssLink,omitempty"`
	PlayStatus string `json:"playStatus,omitempty"`
	Area       string `json:"area,omitempty"`
	Category   string `json:"category,omitempty"`
}

type episode struct {
	SeriesID string `json:"seriesId"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
	Ed2k     string `json:"ed2k,omitempty"`
	Magnet   string `json:"magnet,omitempty"`
}

type service struct {
	domain string
}

func (svc *service) search(query searchQuery) ([]series, error) {
	domain := svc.domain
	url := fmt.Sprintf("http://%s/search/api?keyword=%s&type=resource", domain, query.Keyword)

	r, err := http.Get(url)
	if err != nil {
		log.Printf("无法访问资源服务器: %s, %s\n", domain, err.Error())
		return nil, fmt.Errorf("无法访问资源服务器: %s, %s", domain, err.Error())
	}
	defer r.Body.Close()

	var searchResp struct {
		Data []struct {
			ItemID string `json:"itemid"`
			Title  string `json:"title"`
			Poster string `json:"poster"`
		}
	}
	if err := json.NewDecoder(r.Body).Decode(&searchResp); err != nil {
		log.Printf("无法解析资源服务器返回的数据: %s\n", err.Error())
		return nil, fmt.Errorf("无法解析资源服务器返回的数据: %s", err.Error())
	}

	var ss = make([]series, len(searchResp.Data))

	for i, item := range searchResp.Data {
		s := series{}
		s.ID = item.ItemID
		s.CnName = item.Title
		// convert to large image
		s.Poster = strings.ReplaceAll(item.Poster, "s_", "")
		if query.Details {
			if err := svc.fill(&s); err != nil {
				continue
			}
		}

		ss[i] = s
	}
	return ss, nil
}

func (svc *service) detail(seriesID string) (series, error) {
	s := series{
		ID: seriesID,
	}
	err := svc.fill(&s)
	return s, err
}

func (svc *service) episodes(seriesID string) ([]episode, error) {
	s := series{
		ID: seriesID,
	}
	if err := svc.fill(&s); err != nil {
		log.Printf("无法获取剧集详情: %s, %s\n", seriesID, err.Error())
		return nil, err
	}
	if s.RssLink == "" {
		log.Printf("rssLink not found: %s\n", seriesID)
		return nil, errors.New("rssLink not found")
	}

	r, err := http.Get(s.RssLink)
	if err != nil {
		log.Printf("rssLink 无法访问: %s\n", s.RssLink)
		return nil, fmt.Errorf("rssLink 无法访问: %s", s.RssLink)
	}
	defer r.Body.Close()

	var rssResp struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Items []struct {
				Title  string `xml:"title"`
				Ed2k   string `xml:"ed2k"`
				Magnet string `xml:"magnet"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(r.Body).Decode(&rssResp); err != nil {
		log.Printf("无法解析rss返回的数据: %s, %s\n", s.RssLink, err.Error())
		return nil, fmt.Errorf("无法解析rss返回的数据: %s, %s", s.RssLink, err.Error())
	}

	var episodes = make([]episode, len(rssResp.Channel.Items))

	seasonEpisodeParse := func(name string) (int, int) {
		re := regexp.MustCompile(`(?m)[Ss](\d{1,3})[Ee](\d{1,3})`)
		matches := re.FindStringSubmatch(name)
		if len(matches) != 3 {
			return -1, -1
		}
		season, err := strconv.Atoi(matches[1])
		if err != nil {
			season = -1
		}
		e, err := strconv.Atoi(matches[2])
		if err != nil {
			e = -1
		}
		return season, e
	}

	for i, item := range rssResp.Channel.Items {
		season, e := seasonEpisodeParse(item.Title)
		episodes[i] = episode{
			SeriesID: seriesID,
			Name:     item.Title,
			Season:   season,
			Episode:  e,
			Ed2k:     item.Ed2k,
			Magnet:   item.Magnet,
		}
	}

	return episodes, nil
}

func (svc *service) fill(series *series) error {
	domain := svc.domain
	if series.ID == "" {
		return errors.New("series id must not empty")
	}
	playStatus := make(chan string, 1)
	// get playStatus from api
	go func() {
		url := fmt.Sprintf("http://%s/resource/index_json/rid/%s/channel/tv", domain, series.ID)

		r, err := http.Get(url)
		if err != nil {
			log.Printf("无法获取连载状态: %s", err.Error())
			playStatus <- "无法获取连载状态"
			return
		}
		defer r.Body.Close()

		rawResp, err := ioutil.ReadAll(r.Body)
		if err != nil || len(rawResp) == 0 {
			log.Printf("无法获取连载状态: %s, %s", url, string(rawResp))
			playStatus <- "无法获取连载状态"
			return
		}

		// remove 'var index_info='
		stringResp := string(rawResp)[len("var index_info="):]
		var playStatusResp struct {
			PlayStatus string `json:"play_status"`
		}
		err = json.Unmarshal([]byte(stringResp), &playStatusResp)
		if err != nil {
			log.Printf("无法获取连载状态: %s", err.Error())
			playStatus <- "无法获取连载状态"
			return
		}
		playStatus <- playStatusResp.PlayStatus
	}()

	// get detail from parse page
	url := fmt.Sprintf("http://%s/resource/%s", domain, series.ID)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}
	cnNameHTML, err := doc.Find(".resource-tit h2").Html()
	if err != nil || len(cnNameHTML) == 0 {
		log.Printf("无法获取中文名: %s", cnNameHTML)
		return errors.New("无法获取中文名")
	}

	series.CnName = cnNameHTML[strings.Index(cnNameHTML, "《")+len("《") : strings.Index(cnNameHTML, "》")]
	series.EnName = doc.Find(".resource-con .fl-info li:nth-child(1) > strong").Text()
	series.RssLink = doc.Find(".resource-tit h2 a").AttrOr("href", "")
	series.Area = doc.Find(".resource-con .fl-info li:nth-child(2) > strong").Text()
	series.Category = doc.Find(".resource-con .fl-info li:nth-child(6) > strong").Text()
	if series.Poster == "" {
		series.Poster = doc.Find(".resource-con > div.fl-img > div.imglink > a").AttrOr("href", "")
	}
	series.Link = url
	return nil
}
