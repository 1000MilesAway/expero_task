package main

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

type AggregationJson struct {
	SSCC    string   `json:"sscc"`
	Created string   `json:"created"`
	SGTINs  []string `json:"sgtins"`
}

func (aggr *AggregationJson) validate() (map[string]any, error) {
	if len(aggr.SGTINs) > MaxPacks {
		return gin.H{"ok": false, "error": "Колличество пачек превышает максимум", "error_code": 1}, errValid
	}

	if numeric := regexp.MustCompile(`\d{18}`).MatchString(aggr.SSCC); !numeric {
		return gin.H{"ok": false, "error": "SSCC невалидный", "error_code": 1}, errValid
	}

	gtin := ""
	for _, sgtin := range aggr.SGTINs {
		if len(sgtin) != 27 {
			return gin.H{"ok": false, "error": "SGTIN не 27 символов", "error_code": 1}, errValid
		}
		if numeric := regexp.MustCompile(`\d`).MatchString(sgtin[:14]); !numeric {
			return gin.H{"ok": false, "error": "GTIN не из цифр", "error_code": 1}, errValid
		}
		if abc := regexp.MustCompile(`[\da-zA-Z]+`).MatchString(sgtin[13:]); !abc {
			return gin.H{"ok": false, "error": "GTIN не из букв", "error_code": 1}, errValid
		}

		if gtin != "" {
			if sgtin[:14] != gtin {
				return gin.H{"ok": false, "error": "GTIN отличается", "error_code": 1}, errValid
			}
		}

		gtin = sgtin[:14]
	}

	return nil, nil
}

type DeAggregationJson struct {
	SGTINs []string `json:"sgtins"`
}

func (aggr *DeAggregationJson) validate() (map[string]any, error) {
	for _, sgtin := range aggr.SGTINs {
		if len(sgtin) != 27 {
			return gin.H{"ok": false, "error": "SGTIN не 27 символов", "error_code": 1}, errValid
		}
		if numeric := regexp.MustCompile(`\d`).MatchString(sgtin[:14]); !numeric {
			return gin.H{"ok": false, "error": "GTIN не из цифр", "error_code": 1}, errValid
		}
		if abc := regexp.MustCompile(`[\da-zA-Z]+`).MatchString(sgtin[13:]); !abc {
			return gin.H{"ok": false, "error": "GTIN не из букв", "error_code": 1}, errValid
		}

	}
	return nil, nil
}

type Resp struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
}
