package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

func (b *autoInc) ID_() (id int) {
	b.Lock()
	defer b.Unlock()
	id = b.id
	b.id++
	return
}

type Video struct {
	Id    int
	Titel string
	Url   string
	Dato  string
}

type Person struct {
	Id      int
	Navn    string
	Website string
	Videos  []Video
}

var ai autoInc

func NewPerson() *Person {
	return &Person{
		Id: ai.ID(),
	}
}

func (person *Person) AddVideo(item Video) []Video {
	person.Videos = append(person.Videos, item)
	return person.Videos
}

func main() {
	timeStart := time.Now()
	Videos := []Video{}
	person := Person{
		Id:      NewPerson().Id,
		Navn:    "Micl",
		Website: "test.dk",
		Videos:  Videos,
	}

	c := colly.NewCollector(
		// hvilke Domains er godkendt.
		colly.AllowedDomains("micl-easj.dk"),
		// colly.AllowedDomains("micl-easj.dk", "easj.dk", "easj.cloud.dk", "easj.cloud.panopto.eu", "cloudfront.net"),
	)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		r_validVideoLink, _ := regexp.Compile("(m+p+4)")
		r_findDate, _ := regexp.Compile("([0-9]{2}.[0-9]{2}.[0-9]{4})|([0-9]{4}.[0-9]{2}.[0-9]{2})")
		r_wongDateFormate, _ := regexp.Compile("([0-9]{4}.[0-9]{2}.[0-9]{2})")
		r_otherDatoField, _ := regexp.Compile("/^(?:(?:31(/|-|.)(?:0?[13578]|1[02]))1|(?:(?:29|30)(/|-|.)(?:0?[1,3-9]|1[0-2])2))(?:(?:1[6-9]|[2-9]d)?d{2})$|^(?:29(/|-|.)0?23(?:(?:(?:1[6-9]|[2-9]d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1d|2[0-8])(/|-|.)(?:(?:0?[1-9])|(?:1[0-2]))4(?:(?:1[6-9]|[2-9]d)?d{2})$/g")
		//Bedre Navne.
		mp4LinkCheck := r_validVideoLink.MatchString(link) == true
		copyDataLinkSimple := r_findDate.MatchString(link) == true
		copyLink := r_findDate.Copy().Find([]byte(link))
		usaTimestandCheck := r_wongDateFormate.MatchString(link) == true
		// Han skriver datoen som USA TID YYYY.MM.DD
		// Hvis jeg scanner en anden side der ikke bruger den format, vil den være forkert format.
		copyDataLinkHardCheck := r_otherDatoField.MatchString(link) == true
		titelGotDate := r_findDate.MatchString(e.Text) == true

		if mp4LinkCheck {
			videoTitel := e.Text
			videoDato := ""
			dato := ""
			videoUrl := ("http://micl-easj.Dk" + string(link))
			if copyDataLinkSimple {
				dato = string(copyLink)
				if usaTimestandCheck {
					// USA TID YYYY.MM.DD
					// Laver det om til DD.MM.YYYY
					year := fmt.Sprint(string(dato[:+4]))
					month := fmt.Sprint(string(dato[5:7]))
					date := string(dato[len(dato)-2:])
					videoDato = (date + "." + month + "." + year)
					// fmt.Println("USA TIMESTAND?")
					// fmt.Println("videoUrl: " + videoUrl)
					// fmt.Println("videoDato: " + videoDato)
				} else {
					fmt.Println("Ikke USA TIMESTAND?")
					fmt.Println("videoUrl: " + videoUrl)
					fmt.Println("videoDato: " + videoDato)
					videoDato = string(dato)
				}
			} else if copyDataLinkHardCheck {
				videoDato = string(r_otherDatoField.Copy().Find([]byte(link)))
			} else if titelGotDate {
				videoDato = string(r_otherDatoField.Copy().Find([]byte(e.Text)))
			} else {
				// Kan ikke finde en dato
				videoDato = "N/A"
				fmt.Println("---------------------------")
				fmt.Println("ERROR")
				fmt.Println("---------------------------")
				fmt.Println("DATO ER BLEVET SAT TIL N/A")
				fmt.Println("videoTitel: " + videoTitel)
				fmt.Println("videoUrl: " + videoUrl)
				fmt.Println("videoDato: " + videoDato)
				fmt.Println("---------------------------")
			}
			// Fundet alt data jeg skal bruge. så tilføre det til lidsten.
			newVideo := Video{Id: NewPerson().Id, Titel: videoTitel, Url: videoUrl, Dato: videoDato}
			// Tilføre det til lidsten.
			person.AddVideo(newVideo)
			fmt.Printf("Video_ID: %d -> Dato: %s Titel: %s\n", newVideo.Id, newVideo.Dato, newVideo.Titel)
		} else {
			c.Visit(e.Request.AbsoluteURL(e.Request.ProxyURL + link))
		}
	})
	//c.Visit("https://easj.cloud.panopto.eu/Panopto/Pages/Sessions/List.aspx")
	c.Visit("http://micl-easj.dk/")
	videosJson, _ := json.MarshalIndent(person, "", "")
	f, err := os.OpenFile("videos.json", os.O_APPEND|os.O_WRONLY, 0600)
	// FEJL!
	// Kunne ikke åbne filen.
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// FEJL!
	// Kunne ikke skrive til filen
	if _, err = f.WriteString(string(videosJson)); err != nil {
		panic(err)
	}
	timeEnd := time.Now()
	elapsed := timeEnd.Sub(timeStart)
	//fmt.Print("\033[H\033[2J")
	fmt.Println("---------------------------")
	fmt.Println("DONE DONE")
	fmt.Println("---------------------------")
	fmt.Print("Har tilført ")
	fmt.Print(len(person.Videos))
	fmt.Print(" Videoer til \"./videos.json\"")
	fmt.Println()
	fmt.Print("Det tog hele scanne ")
	fmt.Print(elapsed)
	fmt.Print(" at scanne dem!")
	fmt.Println()
	fmt.Println("---------------------------")

}
