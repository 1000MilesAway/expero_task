package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var MaxPacks = 500

type AggregationJson struct {
	SSCC    string   `json:"sscc"`
	Created string   `json:"created"`
	SGTINs  []string `json:"sgtins"`
}

type DeAggregationJson struct {
	SGTINs []string `json:"sgtins"`
}

type DeAggregationRespJson struct {
	SSCC string `json:"sscc"`
}

type Resp struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
}

func Aggregation(c *gin.Context) {
	var aggr AggregationJson

	err := c.BindJSON(&aggr)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Invalid json data"})
		return
	}
	// logrus.Info(aggr)

	if len(aggr.SGTINs) > MaxPacks {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Колличество пачек превышает максимум", "error_code": 1})
		return
	}

	gtin := aggr.SGTINs[0][:14]
	for i := 1; i < len(aggr.SGTINs); i++ {
		if aggr.SGTINs[i][:14] != gtin {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "GTIN отличается", "error_code": 1})
			return
		}
	}

	ctx := context.Background()

	var dst []interface{}
	err = pgxscan.Select(ctx, GetDB(), &dst, fmt.Sprintf(
		`INSERT INTO packages
		VALUES ('%s', '%s', '{%s}')`, aggr.SSCC, aggr.Created, strings.Join(aggr.SGTINs, ",")))

	if err.Error() == "scanning all: scany: rows final error: ERROR: duplicate key value violates unique constraint \"packages_pkey\" (SQLSTATE 23505)" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "SSCC уже использован", "error_code": 1})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error(), "error_code": 1})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"ok": true})
	return
}

func DeAggregation(c *gin.Context) {
	var aggr DeAggregationJson

	err := c.BindJSON(&aggr)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Invalid json data"})
		return
	}
	logrus.Info(aggr)

	ctx := context.Background()

	result := make(map[string]string)
	for _, sgtin := range aggr.SGTINs {
		var dst []DeAggregationRespJson
		if err := pgxscan.Select(ctx, GetDB(), &dst,
			fmt.Sprintf(`SELECT sscc FROM packages
			WHERE '%s' = ANY(sgtins)`, sgtin)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "SSCC уже использован", "error_code": 1})
			return
		}
		logrus.Info(dst)
		if len(dst) != 0 {
			result[sgtin] = dst[0].SSCC
		} else {
			result[sgtin] = "null"
		}
	}

	c.JSON(http.StatusBadRequest, result)
	return

}
