package PageInfo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"getAcFunPage/utils"
)

// Author: HackerZ
// Time  : 2016-7-5 11:47

// PageInfo Some Infomation of the Page.
type PageInfo struct {
	PageID   string
	Onlooker int64
	Comments int64
	Banana   int64
	// Published string
}

var (
	acfunRoot        = "http://www.acfun.tv"
	acfunContentView = "http://www.acfun.tv/content_view.aspx?"
	ptnPageInfo      = regexp.MustCompile(`\[([0-9]+),([0-9]+),[0-9]+,[0-9]+,[0-9]+,[0-9]+,[0-9]+,([0-9]+)\]`)
	// ptnPageInfo      = regexp.MustCompile(`<span class="pts">([0-9]+)</span><span>围观</span>&nbsp;·&nbsp;&nbsp;<span class="pts pointer">([0-9]+)</span><span class="pinglun">评论</span>&nbsp;·&nbsp;&nbsp;<span class="pts">([0-9]+)</span><span>香蕉</span>&nbsp;/&nbsp;&nbsp;<span class="time">发布于 ([0-9]+年 [0-9]+月[0-9]+日 [0-9]+:[0-9]+)</span>`)
)

// GetPageInfo Return Page Info which pageID.
func GetPageInfo(pageID string) PageInfo {
	// Replace the string "/a/ac"
	pageID = strings.Replace(pageID, "/a/ac", "", 1)

	// GetPageInfo By GET Request.
	raw, statusCode := getPageInfo(acfunContentView + "contentId=" + pageID + "&channelId=110")
	if statusCode != "200 OK" {
		fmt.Printf("Get %s PageInfo Error.\n", pageID)
		return PageInfo{
			PageID:   pageID,
			Onlooker: -1,
			Banana:   -1,
			Comments: -1,
		}
	}

	// Matching PageInfo.
	pageInfo := findPageInfo(pageID, raw)
	return pageInfo
}

func getPageInfo(url string) (content, statusCode string) {
	resp, err := http.Get(url)

	if err != nil {
		statusCode = "-1"
		fmt.Println("Get PageInfo Error -->", err.Error())
		return
	}
	defer resp.Body.Close()

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		statusCode = "-2"
		fmt.Println("Get PageInfo Error -->", err.Error())
		return
	}
	statusCode = resp.Status
	content = string(bys)
	return
}

func findPageInfo(pageID, content string) PageInfo {
	matches := ptnPageInfo.FindStringSubmatch(content)

	// fmt.Println("PageInfo", matches)
	// Request: http://www.acfun.tv/content_view.aspx?contentId=2867906&channelId=110
	// Output : [2085,72,0,0,0,38,1,1]
	return PageInfo{
		pageID,
		utils.StrToInt64(matches[1]),
		utils.StrToInt64(matches[2]),
		utils.StrToInt64(matches[3]),
	}

}
