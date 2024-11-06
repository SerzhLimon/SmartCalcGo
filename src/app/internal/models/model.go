package models

import "gopkg.in/guregu/null.v4"

type PushNumberRequest struct {
	Value int `json:"value"`
}

type PushOperatorRequest struct {
	Value string `json:"value"`
}

type OutputResponse struct {
	Output string      `json:"result"`
	Error  null.String `json:"errors"`
}

type MakeResultRequest struct {
	Input string `json:"result"`
}

type GraphRequest struct {
	Input string `json:"result"`
	Xmin  string `json:"xmin"`
	Xmax  string `json:"xmax"`
}

type GraphResponse struct {
	Image string      `json:"image"`
	Error null.String `json:"errors"`
}

type HistoryResponse struct {
	History [][]string `json:"history"`
}
