package controller

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart"
	"net/http"
	"star-trend/github"
	"time"
)

type repoParam struct {
	Owner string `uri:"owner" binding:"required"`
	Repo  string `uri:"repo" binding:"required"`
}

func RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Repo": github.Repository{},
	})
}

func RenderGraph(c *gin.Context) {
	var r repoParam
	if err := c.ShouldBindUri(&r); err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": err.Error(),
		})
		return
	}

	repo, err := github.Gh.GetRepoDetail(fmt.Sprintf("%s/%s", r.Owner, r.Repo))
	if err != nil {
		log.WithError(err).Error("gh.GetRepoDetail failed")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": fmt.Sprintf("GetRepoDetail failed: %v", err.Error()),
		})
		return
	}

	stargazers, err := github.Gh.GetStargazers(repo)
	if err != nil {
		log.WithError(err).Error("gh.GetStargazers failed")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": fmt.Sprintf("GetStargazers failed: %v", err.Error()),
		})
		return
	}

	if len(stargazers) == 0 {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": "",
		})
		return
	}

	XValues := make([]time.Time, 0)
	YValues := make([]float64, 0)
	for i, v := range stargazers {
		XValues = append(XValues, v.StarredAt)
		YValues = append(YValues, float64(i+1))
	}

	graph := chart.Chart{
		Title:      fmt.Sprintf("%s Star Trend", repo.FullName),
		TitleStyle: chart.StyleShow(),
		XAxis: chart.XAxis{
			Name:           "Date",
			NameStyle:      chart.StyleShow(),
			ValueFormatter: chart.TimeDateValueFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name:           "Star Count",
			NameStyle:      chart.StyleShow(),
			ValueFormatter: chart.FloatValueFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name: "Star Count",
				Style: chart.Style{
					//StrokeWidth: chart.Disabled,
					//DotWidth:    5,
					//DotColor: chart.ColorBlue,
					Show: true,
				},
				//XValues: chart.TimeSeriesValues(starData, func(i int) time.Time {
				//	return starData[i].Date
				//}),
				//YValues: chart.TimeSeriesValues(starData, func(i int) float64 {
				//	return float64(starData[i].StarCount)
				//}),
				XValues: XValues,
				YValues: YValues,
			},
		},
	}

	c.Header("Content-Type", "image/png")
	graph.Render(chart.PNG, c.Writer)
}

func RenderRepoDetail(c *gin.Context) {
	fullName := c.PostForm("repository")
	if fullName == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": "Please enter owner/repo for example: frankxjkuang/star-trend",
		})
		return
	}

	repo, err := github.Gh.GetRepoDetail(fullName)
	if err != nil {
		log.WithError(err).Error("gh.GetRepoDetail failed")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Repo":     github.Repository{},
			"ErrorMsg": fmt.Sprintf("GetRepoDetail failed: %v", err.Error()),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Repo":     repo,
		"ErrorMsg": "",
	})
}
