package incruit

import (
	"fmt"
	"main/package/helper"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type InputData struct {
	Url   string
	Count string
}

type ScrapeJob struct {
	Status         bool
	Link           string
	Title          string
	Condition      string
	Corp           string
	RegisteredDate string
	Sequence       int
}

type Pagenation struct {
	Total           int
	TotalPage       int
	IsShowNextPage  bool
	IsShowPrevPage  bool
	IsShowPageCount bool
	NextPage        string
	PrevPage        string
}

// 사이트 주소와 가져올 갯수를 입력 받는다.
func New(url, count string) *InputData {
	return &InputData{url, count}
}

// 스크랩
func (i InputData) GetData(target, keyword, page string) ([]ScrapeJob, Pagenation) {
	return getPage(target, keyword, page, i.Url, i.Count)
}

// 페이지의 정보를 가져온다
func getPage(target, keyword, page, siteHome, count string) ([]ScrapeJob, Pagenation) {
	searchUrl := getUrl(keyword, page, siteHome)
	jobs := []ScrapeJob{}
	pagenation := Pagenation{}
	secondCh := make(chan ScrapeJob)

	resp, err := http.Get(searchUrl)
	helper.CheckErr(err)
	helper.CheckConnect(resp)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	helper.CheckErr(err)
	fmt.Println("incruit Connected: ", resp.Status)

	jobItems := doc.Find(".section .litype01 li")

	if jobItems.Length() > 0 {
		jobItems.Each(func(i int, s *goquery.Selection) {
			go getExtractJob(s, searchUrl, siteHome, secondCh)
		})

		for i := 0; i < jobItems.Length(); i++ {
			job := <-secondCh
			jobs = append(jobs, job)
		}

		// 등록일순으로, 등록일이 같으면 제목순으로 정렬한다.
		sort.SliceStable(jobs, func(i, j int) bool {
			if jobs[i].Sequence == jobs[j].Sequence {
				return jobs[i].Title < jobs[j].Title
			}
			return jobs[i].Sequence < jobs[j].Sequence
		})

		// 검색결과의 페이지네이션 정보를 정리한다.
		totalResult := doc.Find(".section .numall").Text()

		total, errAtoi := strconv.Atoi(strings.Join(helper.GetNumFromString(totalResult), ""))
		helper.CheckErr(errAtoi)

		count, errAtoi := strconv.Atoi(count)
		helper.CheckErr(errAtoi)

		totalPage := total / count
		if totalPage >= 1 {
			rest := total % count
			if rest != 0 {
				totalPage++
			}
		}

		isShowNextPage := true
		isShowPrevPage := false

		currentPage, errAtoi := strconv.Atoi(page)
		helper.CheckErr(errAtoi)

		next := map[string]string{"target": target, "keyword": keyword, "page": strconv.Itoa(currentPage * count)}
		nextPage := helper.UrlParamBuild(next, helper.UrlDirectoryBuild([]string{"search"}, ""))
		prevPage := ""

		if currentPage >= totalPage {
			isShowNextPage = false
		}

		if currentPage > 1 {
			isShowPrevPage = true
			prev := map[string]string{"target": target, "keyword": keyword, "page": strconv.Itoa((currentPage * count) - count)}
			prevPage = helper.UrlParamBuild(prev, helper.UrlDirectoryBuild([]string{"search"}, ""))
		}

		isShowPageCount := false
		if isShowNextPage || isShowPrevPage {
			isShowPageCount = true
		}

		pagenation = Pagenation{
			Total:           total,
			TotalPage:       totalPage,
			IsShowNextPage:  isShowNextPage,
			IsShowPrevPage:  isShowPrevPage,
			IsShowPageCount: isShowPageCount,
			NextPage:        nextPage,
			PrevPage:        prevPage,
		}
	} else {
		noResults := ScrapeJob{
			Status: false,
		}
		jobs = append(jobs, noResults)
	}

	return jobs, pagenation
}

// 구직 정보를 가져온다.
func getExtractJob(s *goquery.Selection, searchUrl, siteHome string, secondCh chan<- ScrapeJob) {
	href, errFind := s.Find(".section .litype01 li .rcrtTitle a").Attr("href")
	if !errFind {
		href = "#"
	}
	title := s.Find(".section .litype01 li .rcrtTitle a").Text()
	if title == "" {
		title = "제목없음"
	}
	condition := helper.CleanString(s.Find(".section .litype01 li .etc span").Text())
	corp := helper.CleanString(s.Find(".section .litype01 li h3 a").Text())
	registeredDate := s.Find(".section .litype01 li .info").Text()
	if registeredDate == "" {
		registeredDate = "00"
	}

	intDate := helper.GetNumFromString(registeredDate)

	fmt.Println(href, title, condition, corp, registeredDate, intDate)
	// sequence, errAtoi := strconv.Atoi(strings.Join(intDate, ""))
	// helper.CheckErr(errAtoi)

	secondCh <- ScrapeJob{
		Status:         true,
		Link:           siteHome + href,
		Title:          title,
		Condition:      condition,
		Corp:           corp,
		RegisteredDate: registeredDate,
		// Sequence:       sequence,
	}
	fmt.Println(secondCh)
}

// 필요한 url을 반환한다.
func getUrl(keyword, page, url string) string {
	var siteUrl string
	var params map[string]string

	directorys := []string{"list", "search.asp"}
	siteUrl = helper.UrlDirectoryBuild(directorys, url)

	params = map[string]string{
		"col":     "job",
		"src":     "gsw*search",
		"kw":      keyword,
		"startno": page,
	}

	return helper.UrlParamBuild(params, siteUrl)
}

// https://www.incruit.co.kr/Search/?stext=golang&tabType=recruit&Page_No=1

// 전체 페이지네이션 갯수를 반환한다.
// func getTotalPage(url, attr string) int {
// 	var countTotalPage int

// 	resp, err := http.Get(url)
// 	helper.CheckErr(err)
// 	helper.CheckConnect(resp)
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	helper.CheckErr(err)

// 	doc.Find(attr).Each(func(i int, s *goquery.Selection) {
// 		countTotalPage = s.Find("a").Length()
// 	})

// 	return countTotalPage
// }
