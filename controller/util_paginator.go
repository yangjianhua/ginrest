package controller

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Package from https://github.com/biezhi/gorm-paginator
// Since it's a single file package, copy to here

type PagingParam struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

type Paginator struct {
	TotalCount int         `json:"total_count"`
	TotalPage  int         `json:"total_page"`
	Data       interface{} `json:"data"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	PrevPage   int         `json:"prev_page"`
	NextPage   int         `json:"next_page"`
}

func Pagging(p *PagingParam, dataSource interface{}) *Paginator {
	db := p.DB
	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int
	var offset int
	go countRecords(db, dataSource, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}
	db.Offset(offset).Limit(p.Limit).Find(dataSource)
	<-done

	paginator.TotalCount = count
	paginator.Data = dataSource
	paginator.Page = p.Page
	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, countDataSource interface{}, done chan bool, count *int) {
	db.Model(countDataSource).Count(count)
	done <- true
}
