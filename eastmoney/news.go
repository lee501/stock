package eastmoney

import (
	"encoding/json"
	"fmt"
)

// ════════════════════════════════════════
// 个股新闻
// ════════════════════════════════════════

type StockNews struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Source  string `json:"source"`
	URL     string `json:"url"`
	Summary string `json:"summary,omitempty"`
}

func GetStockNews(code string, limit int) ([]StockNews, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	u := fmt.Sprintf(
		"https://search-api-web.eastmoney.com/search/jsonp?cb=&param=%%7B%%22uid%%22:%%22%%22,%%22keyword%%22:%%22%s%%22,%%22type%%22:%%5B%%22cmsArticleWebOld%%22%%5D,%%22client%%22:%%22web%%22,%%22clientType%%22:%%22web%%22,%%22clientVersion%%22:%%22curr%%22,%%22param%%22:%%7B%%22cmsArticleWebOld%%22:%%7B%%22searchScope%%22:%%22default%%22,%%22sort%%22:%%22default%%22,%%22pageIndex%%22:1,%%22pageSize%%22:%d,%%22preTag%%22:%%22%%22,%%22postTag%%22:%%22%%22%%7D%%7D%%7D",
		code, limit,
	)
	body, err := DoGet(u)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Result []struct {
			CmsArticleWebOld []struct {
				Title      string `json:"title"`
				Date       string `json:"date"`
				MediaName  string `json:"mediaName"`
				ArticleURL string `json:"url"`
				Content    string `json:"content"`
			} `json:"cmsArticleWebOld"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse news failed: %v", err)
	}

	var news []StockNews
	if len(raw.Result) > 0 {
		for _, item := range raw.Result[0].CmsArticleWebOld {
			n := StockNews{
				Title:  item.Title,
				Date:   item.Date,
				Source: item.MediaName,
				URL:    item.ArticleURL,
			}
			if len(item.Content) > 200 {
				n.Summary = item.Content[:200] + "..."
			} else {
				n.Summary = item.Content
			}
			news = append(news, n)
		}
	}
	return news, nil
}
