package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
)

const (
	artist = "bekar"
)

type (
	provider interface {
		scrap()
		getEndpoint() string
		getName() string
	}

	storage struct {
		*sync.Mutex
		values map[string]*exportedProvider
	}

	exportedProvider struct {
		Date    string            `yaml:"date"`
		Address string            `yaml:"address"`
		Tickets map[string]string `yaml:"tickets"`
	}

	tmProvider struct{}
	stProvider struct{}
)

func (s *storage) store(e exportedProvider, p provider) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.values[e.Date] == nil {
		s.values[e.Date] = &e
	} else {
		s.values[e.Date].Tickets[p.getName()] = e.Tickets[p.getName()]
	}
}

func (t *tmProvider) scrap() {
	c := colly.NewCollector()
	c.OnHTML("#resultsListZone", func(e *colly.HTMLElement) {
		e.DOM.Children().Find(".bloc-result-content").Each(func(_ int, s *goquery.Selection) {
			url := s.Find("#urlToConcertHallLabel")
			e := exportedProvider{
				Tickets: make(map[string]string),
			}
			e.Address = url.Text()
			e.Tickets[t.getName()], _ = url.Attr("href")

			d, _ := s.Find("time").Attr("content")

			e.Date = d

			store.store(e, t)
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(t.getEndpoint())
}

func (t *tmProvider) getEndpoint() string {
	return "https://www.ticketmaster.fr/fr/resultat?ipSearch=" + artist
}

func (t *tmProvider) getName() string {
	return "ticketmaster"
}

func (st *stProvider) scrap() {
	c := colly.NewCollector()
	c.OnHTML("#search-results-wrapper", func(e *colly.HTMLElement) {
		e.DOM.Children().Find(".g-blocklist-link").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			e := exportedProvider{
				Tickets: make(map[string]string),
			}
			e.Tickets[st.getName()] = "https://www.seetickets.com" + href
			values := strings.Split(s.Find(".g-blocklist-sub-text").Text(), "\n")
			e.Address = strings.TrimSpace(values[4])

			d := ""
			s.Find("time").Each(func(x int, se *goquery.Selection) {
				if x != 1 {
					return
				}

				d, _ = se.Attr("datetime")
			})

			date, _ := time.Parse("02 01 2006", d)

			e.Date = date.Format("2006-01-02")

			store.store(e, st)
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(st.getEndpoint())
}

func (s *stProvider) getEndpoint() string {
	return "https://www.seetickets.com/fr/search?q=" + artist
}

func (s *stProvider) getName() string {
	return "seetickets"
}

var (
	providers = []provider{
		&tmProvider{},
		&stProvider{},
	}
	store = &storage{
		Mutex:  &sync.Mutex{},
		values: make(map[string]*exportedProvider),
	}
)

func main() {
	var wg sync.WaitGroup

	for _, pr := range providers {
		wg.Add(1)
		go func(wgrp *sync.WaitGroup, p provider) {
			defer wg.Done()
			p.scrap()
		}(&wg, pr)
	}

	wg.Wait()

	vs := make([]exportedProvider, 0)
	for _, v := range store.values {
		vs = append(vs, *v)
	}
	b, _ := yaml.Marshal(vs)
	ioutil.WriteFile("../../data/bekar/dates.yaml", b, 0755)
}
