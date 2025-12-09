package helpers

import (
	"net/http"
	"slices"
	"strings"

	"beardsall.xyz/golanghttpplayground/config"
	"beardsall.xyz/golanghttpplayground/repository"
)

func ExtractQueryParams(req *http.Request) ([]repository.QueryFilter, string) {
	// func ExtractQueryParams(req *http.Request) (map[string]any, string) {
	// queryParams := make(map[string]any)
	queryOperators := make(map[string]string)
	FullQueryParams := []repository.QueryFilter{}
	linkOperator := "AND"

	for param, value := range req.URL.Query() {
		if slices.Contains(config.RESERVED_PARAMS, param) {
			continue
		}
		if strings.HasSuffix(param, config.OPERATOR_SUFFIX) {
			queryOperators[param] = value[0]
			continue
		}
		if param == config.LINK_OPERATOR_PARAM {
			linkOperator = strings.ToUpper(value[0])
			continue
		}
		FullQueryParams = append(FullQueryParams, repository.QueryFilter{
			FieldName: param,
			Operator:  config.DEFAULT_OPERATOR,
			Value:     value[0],
		})
		// queryParams[param] = value[0]
	}

	// for paramIndex := range FullQueryParams {
	// 	// if the param is in queryOperators, set its operator, else set its fallback to default (config.DEFAULT_OPERATOR)
	// keys := make([]string, 0, len(queryOperators))
	// for k := range queryOperators {
	// 	keys = append(keys, k)
	// }
	// 	if slices.Contains(keys, FullQueryParams[paramIndex].FieldName+"Operator") {
	// 		FullQueryParams[paramIndex].Operator = queryOperators[FullQueryParams[paramIndex].FieldName+"Operator"]
	// 	}
	// }

	filters := make([]string, 0, len(queryOperators))
	for fqpIndex := range FullQueryParams {
		filters = append(filters, FullQueryParams[fqpIndex].FieldName)
	}

	for key, operator := range queryOperators {
		print(key)
		op, _ := strings.CutSuffix(key, config.OPERATOR_SUFFIX)
		if !slices.Contains(filters, op) {
			continue
		}
		queryParamIndex := indexOf(op, FullQueryParams)
		FullQueryParams[queryParamIndex].Operator = operator
	}

	// TODO: make the and dynamic based on a query param, but if operator is provided instead of <field>Operator it should use that, also make it so you can use that as a default and override certain fields
	return FullQueryParams, linkOperator
	// return FullQueryParams, "AND"
	// return queryParams, "AND"
}

func indexOf(fieldName string, data []repository.QueryFilter) int {
	for k, v := range data {
		if fieldName == v.FieldName {
			return k
		}
	}
	return -1
}
