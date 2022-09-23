package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
	"time"
)

type Auto struct {
	Brand string
	Age   string
	Cost  string

	Mileage      string
	Transmission string
	EnginePower  string
	EnginesType  string
	DriveUnit    string

	DromLink string
	IsSell   string
}

const Filename = "result.json"

func main() {
	var newFileName string = "productsNew.json"
	// Instantiate default collector
	c := colly.NewCollector()
	autos := make([]Auto, 0)
	auto := Auto{}
	//flats := make([]Flat, 0)
	c.OnHTML("a.css-5l099z.ewrty961", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Println(link)
		auto.DromLink = link
		e.ForEach("div.css-13ocj84.e727yh30", func(i int, z *colly.HTMLElement) {
			z.ForEach("div", func(i int, u *colly.HTMLElement) {
				if u.ChildText("div.css-r91w5p.e3f4v4l2") == "" {
					u.ForEach("div.css-17lk78h.e3f4v4l2", func(i int, v *colly.HTMLElement) {
						title := v.ChildText("span")
						nameAndDate := strings.Split(title, ", ")
						auto.Brand = nameAndDate[0]
						auto.Age = nameAndDate[1]
						auto.IsSell = "Продаётся"
					})
				} else {
					u.ForEach("div.css-r91w5p.e3f4v4l2", func(i int, v *colly.HTMLElement) {
						title := v.ChildText("span")
						nameAndDate := strings.Split(title, ", ")
						auto.Brand = nameAndDate[0]
						auto.Age = nameAndDate[1]
						auto.IsSell = "Снят с прожажи"
					})
				}
			})
			z.ForEach("div.css-1fe6w6s.e162wx9x0", func(i int, v *colly.HTMLElement) {
				description := v.ChildText("span.css-1l9tp44.e162wx9x0")
				//description := v.ChildText("span.css-1l9tp44.e162wx9x0")
				//description = strings.ReplaceAll(description, " ", "")
				descriptionArr := strings.Split(description, ",")

				for i := range descriptionArr {
					if strings.Contains(descriptionArr[i], "л.с") {
						auto.EnginePower = descriptionArr[i]
					} else if strings.Contains(descriptionArr[i], "АКПП") || strings.Contains(descriptionArr[i], "механика") || strings.Contains(descriptionArr[i], "робот") || strings.Contains(descriptionArr[i], "вариатор") {
						auto.Transmission = descriptionArr[i]
					} else if strings.Contains(descriptionArr[i], "бензин") || strings.Contains(descriptionArr[i], "дизель") || strings.Contains(descriptionArr[i], "гибрид") {
						auto.DriveUnit = descriptionArr[i]
					} else if strings.Contains(descriptionArr[i], "передний") || strings.Contains(descriptionArr[i], "задний") || strings.Contains(descriptionArr[i], "4DW") {
						auto.EnginesType = descriptionArr[i]
					} else if strings.Contains(descriptionArr[i], "тыс. км") {
						descriptionArr[i] = strings.ReplaceAll(descriptionArr[i], "\u003c", "менее")
						auto.Mileage = descriptionArr[i]
					}

				}
			})
		})
		e.ForEach("div.css-1dkhqyq.ep0qbyc0", func(i int, z *colly.HTMLElement) {
			z.ForEach("div", func(i int, u *colly.HTMLElement) {
				/*
					u.ForEach("div.css-1i8tk3y.eyvqki92", func(i int, v *colly.HTMLElement) {
						status := v.Text
						if strings.Contains(status, "снят с продажи") {
							auto.IsSell = status
						} else {
							auto.IsSell = "Ещё не продан"
						}
					})

				*/
				u.ForEach("div.css-1i8tk3y.eyvqki92", func(i int, v *colly.HTMLElement) {
					v.ForEach("div.css-1dv8s3l.eyvqki91", func(i int, p *colly.HTMLElement) {
						p.ForEach("span.css-46itwz.e162wx9x0", func(i int, r *colly.HTMLElement) {
							costs := r.Text
							costs = strings.ReplaceAll(costs, "\u00a0", "")
							//fmt.Println(costs)
							auto.Cost = costs
							//fmt.Println(costs)
						})
					})
				})
			})
		})

		autos = append(autos, auto)
	})

	c.OnHTML("a.css-4gbnjj.e24vrp30", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(autos, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err = os.WriteFile(newFileName, js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})
	numVisited := 0
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		if numVisited > 50 {
			r.Abort()
		}
		numVisited++
	})

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: (1 * time.Second) / 2,
		// Add an additional random delay
		//RandomDelay: 1 * time.Second,
	})

	c.Visit("https://auto.drom.ru/region70/")
}
