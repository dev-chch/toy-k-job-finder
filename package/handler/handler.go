package handler

import (
	"fmt"
	"main/package/helper"
	"main/package/scrape/incruit"
	"main/package/scrape/indeed"
	"main/package/scrape/jobkorea"
	"main/package/scrape/saramin"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// 첫 페이지
func Index(e echo.Context) error {
	fmt.Println("Index Page")
	return e.Render(http.StatusOK, "index.html", map[string]interface{}{
		"title": "구직 통합검색 토이 프로젝트",
	})
}

// 검색 페이지
func Search(e echo.Context) error {
	target := e.QueryParam("target")
	if target == "" {
		target = "saramin"
	}

	keyword := e.QueryParam("keyword")
	page := e.QueryParam("page")
	if page == "" {
		page = "1"
	}
	fmt.Println("Search Page Target :", target, ", Keyword :", keyword, ", Page :", page)

	// 검색어가 없으면 검색 결과가 없습니다. 페이지를 표시한다.
	if keyword == "" {
		return e.Render(http.StatusOK, "search.html", map[string]interface{}{
			"status":  false,
			"title":   "검색 결과가 없습니다.",
			"keyword": keyword,
		})
	}

	switch target {
	case "saramin":
		saramin := saramin.New("https://www.saramin.co.kr", "40")
		jobsList, pagenation := saramin.GetData(target, keyword, page)

		// 검색 결과가 없으면 검색 결과가 없습니다. 페이지를 표시한다.
		for _, job := range jobsList {
			if !job.Status {
				return e.Render(http.StatusOK, "search.html", map[string]interface{}{
					"status":  false,
					"title":   "검색 결과가 없습니다.",
					"keyword": keyword,
					"target":  target,
				})
			}
		}

		total := helper.NumberFormat(int64(pagenation.Total))
		return e.Render(http.StatusOK, "search.html", map[string]interface{}{
			"status":          true,
			"title":           "'" + keyword + "' 검색 결과",
			"keyword":         keyword,
			"jobsList":        jobsList,
			"jobsTotal":       total,
			"currentPage":     page,
			"totalPage":       pagenation.TotalPage,
			"isShowNextPage":  pagenation.IsShowNextPage,
			"isShowPrevPage":  pagenation.IsShowPrevPage,
			"isShowPageCount": pagenation.IsShowPageCount,
			"nextPage":        pagenation.NextPage,
			"prevPage":        pagenation.PrevPage,
			"target":          target,
		})
	case "jobkorea":
		jobkorea := jobkorea.New("https://www.jobkorea.co.kr", "40")
		jobsList, pagenation := jobkorea.GetData(target, keyword, page)

		// 검색 결과가 없으면 검색 결과가 없습니다. 페이지를 표시한다.
		for _, job := range jobsList {
			if !job.Status {
				return e.Render(http.StatusOK, "search.html", map[string]interface{}{
					"status":  false,
					"title":   "검색 결과가 없습니다.",
					"keyword": keyword,
					"target":  target,
				})
			}
		}

		total := helper.NumberFormat(int64(pagenation.Total))
		return e.Render(http.StatusOK, "search.html", map[string]interface{}{
			"status":          true,
			"title":           "'" + keyword + "' 검색 결과",
			"keyword":         keyword,
			"jobsList":        jobsList,
			"jobsTotal":       total,
			"currentPage":     page,
			"totalPage":       pagenation.TotalPage,
			"isShowNextPage":  pagenation.IsShowNextPage,
			"isShowPrevPage":  pagenation.IsShowPrevPage,
			"isShowPageCount": pagenation.IsShowPageCount,
			"nextPage":        pagenation.NextPage,
			"prevPage":        pagenation.PrevPage,
			"target":          target,
		})
	case "indeed":
		indeed := indeed.New("https://kr.indeed.com", "10")
		jobsList, pagenation := indeed.GetData(target, keyword, page)

		// 검색 결과가 없으면 검색 결과가 없습니다. 페이지를 표시한다.
		for _, job := range jobsList {
			if !job.Status {
				return e.Render(http.StatusOK, "search.html", map[string]interface{}{
					"status":  false,
					"title":   "검색 결과가 없습니다.",
					"keyword": keyword,
				})
			}
		}

		total := helper.NumberFormat(int64(pagenation.Total))
		return e.Render(http.StatusOK, "search.html", map[string]interface{}{
			"status":          true,
			"title":           "'" + keyword + "' 검색 결과",
			"keyword":         keyword,
			"jobsList":        jobsList,
			"jobsTotal":       total,
			"currentPage":     page,
			"totalPage":       pagenation.TotalPage,
			"isShowNextPage":  pagenation.IsShowNextPage,
			"isShowPrevPage":  pagenation.IsShowPrevPage,
			"isShowPageCount": pagenation.IsShowPageCount,
			"nextPage":        pagenation.NextPage,
			"prevPage":        pagenation.PrevPage,
			"target":          target,
		})
	case "incruit":
		incruit := incruit.New("https://search.incruit.com/", "40")
		jobsList, pagenation := incruit.GetData(target, keyword, page)

		// 검색 결과가 없으면 검색 결과가 없습니다. 페이지를 표시한다.
		for _, job := range jobsList {
			if !job.Status {
				return e.Render(http.StatusOK, "search.html", map[string]interface{}{
					"status":  false,
					"title":   "검색 결과가 없습니다.",
					"keyword": keyword,
				})
			}
		}

		total := helper.NumberFormat(int64(pagenation.Total))
		return e.Render(http.StatusOK, "search.html", map[string]interface{}{
			"status":          true,
			"title":           "'" + keyword + "' 검색 결과",
			"keyword":         keyword,
			"jobsList":        jobsList,
			"jobsTotal":       total,
			"currentPage":     page,
			"totalPage":       pagenation.TotalPage,
			"isShowNextPage":  pagenation.IsShowNextPage,
			"isShowPrevPage":  pagenation.IsShowPrevPage,
			"isShowPageCount": pagenation.IsShowPageCount,
			"nextPage":        pagenation.NextPage,
			"prevPage":        pagenation.PrevPage,
			"target":          target,
		})
	}

	return e.Render(http.StatusOK, "search.html", map[string]interface{}{
		"status":  false,
		"title":   "검색 결과가 없습니다.",
		"keyword": "",
	})
}

func CurrentURL(r *http.Request) string {
	hostname, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	return hostname + r.URL.Path
}

// 소개 페이지
func About(c echo.Context) error {
	fmt.Println("About Page")
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"title": "ABOUT",
	})
}
