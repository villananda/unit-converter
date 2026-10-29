package main

import (
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// section is the view model for one category's form and result.
type section struct {
	Key       string
	Label     string
	Units     []Unit
	Value     string
	From      string
	To        string
	FromLabel string
	ToLabel   string
	Result    string
	Error     string
}

func defaultSection(c Category) section {
	return section{
		Key:   c.Key,
		Label: c.Label,
		Units: c.Units,
		Value: "1",
		From:  c.Units[0].Key,
		To:    c.Units[1].Key,
	}
}

func allSections() []section {
	sections := make([]section, 0, len(categories))
	for _, c := range categories {
		sections = append(sections, defaultSection(c))
	}
	return sections
}

func formatResult(v float64) string {
	return strconv.FormatFloat(v, 'g', 10, 64)
}

// renderPage renders the tabbed page with the given category selected.
func renderPage(c *gin.Context, status int, sections []section, active string) {
	c.HTML(status, "index.html", gin.H{"Sections": sections, "Active": active})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		active := categories[0].Key
		if cat, ok := findCategory(c.Query("tab")); ok {
			active = cat.Key
		}
		renderPage(c, http.StatusOK, allSections(), active)
	})

	r.POST("/convert/:category", func(c *gin.Context) {
		cat, ok := findCategory(c.Param("category"))
		if !ok {
			c.String(http.StatusNotFound, "category not found")
			return
		}

		sections := allSections()
		for i := range sections {
			if sections[i].Key == cat.Key {
				sections[i] = convertSection(c, cat)
			}
		}
		renderPage(c, http.StatusOK, sections, cat.Key)
	})

	log.Fatal(r.Run(":8080"))
}

// convertSection reads the submitted form for cat and fills in the result or error.
func convertSection(c *gin.Context, cat Category) section {
	s := defaultSection(cat)
	s.Value = strings.TrimSpace(c.PostForm("value"))
	s.From = c.PostForm("from")
	s.To = c.PostForm("to")

	value, err := strconv.ParseFloat(s.Value, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		s.Error = "Enter a valid number."
		return s
	}
	result, err := Convert(cat.Key, s.From, s.To, value)
	if err != nil {
		if errors.Is(err, ErrBelowZero) {
			s.Error = err.Error()
		} else {
			s.Error = "Invalid unit."
		}
		return s
	}
	src, _ := cat.unit(s.From)
	dst, _ := cat.unit(s.To)
	s.FromLabel = src.Label
	s.ToLabel = dst.Label
	s.Result = formatResult(result)
	return s
}
