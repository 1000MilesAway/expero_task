package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var MaxPacks = 500
var errValid = errors.New("Validation error")

func Aggregation(c *gin.Context) {
	var aggr AggregationJson

	err := c.BindJSON(&aggr)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Invalid json data"})
		return
	}

	resp, err := aggr.validate()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err = GetDB().Query("INSERT INTO packages VALUES ($1, $2, $3)", aggr.SSCC, aggr.Created, pq.Array(aggr.SGTINs))

	if sqlError, ok := err.(*pq.Error); ok {
		if sqlError.Code == "23505" {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "SSCC уже использован", "error_code": 1})
			return
		}
	} else if err != nil {
		logrus.Error(err)
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
	resp, err := aggr.validate()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	result := make(map[string]string)
	for _, sgtin := range aggr.SGTINs {
		rows, err := GetDB().Query("SELECT sscc FROM packages WHERE $1 = ANY(sgtins)", sgtin)
		defer rows.Close()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "SSCC уже использован", "error_code": 1})
			return
		}

		var ssccList []string
		for rows.Next() {
			var temp string
			rows.Scan(&temp)
			ssccList = append(ssccList, temp)
		}
		if len(ssccList) != 0 {
			result[sgtin] = ssccList[0]
		} else {
			result[sgtin] = "null"
		}
	}

	c.JSON(http.StatusBadRequest, result)
	return

}
