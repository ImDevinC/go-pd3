package models

import "time"

type PD3Response struct {
	Data   []PD3DataResponse `json:"data"`
	Paging PD3PagingResponse `json:"paging"`
}

type PD3DataResponse struct {
	RecordId  string               `json:"recordId"`
	Namespace string               `json:"namespace"`
	UserId    string               `json:"userId"`
	Challenge PD3ChallengeResponse `json:"challenge"`
	Progress  PD3ProgressResponse  `json:"progress"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Status    string               `json:"status"`
	IsActive  bool                 `json:"isActive"`
}

type PD3ChallengeResponse struct {
	ChallengeId  string                  `json:"challengeId"`
	Namespace    string                  `json:"namespace"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Prerequisite PD3PrerequisiteResponse `json:"prerequisite"`
	Objective    PD3ObjectiveResponse    `json:"objective"`
	Reward       PD3RewardResponse       `json:"reward"`
	Tags         []string                `json:"tags"`
	OrderNo      int                     `json:"orderNo"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
	IsActive     bool                    `json:"isActive"`
}

type PD3PrerequisiteResponse struct {
	Stats                 []PD3StatsResponse `json:"stats"`
	Items                 interface{}        `json:"items"`                 // All of these values are currently empty, will revisit later
	CompletedChallengeIds interface{}        `json:"completedChallengeIds"` // This can either be a list of strings, or a list of objects in the form of {challengeId: string, isCompleted: bool}
}

type PD3ObjectiveResponse struct {
	Stats []PD3StatsResponse `json:"stats"`
}

type PD3StatsResponse struct {
	StatCode string `json:"statCode"`
	Value    int    `json:"value"`
}

type PD3RewardResponse struct {
	RewardId string             `json:"rewardId"`
	Stats    []PD3StatsResponse `json:"stats"`
	Items    interface{}        `json:"items"` // All of these values are currently empty, will revisit later
}

type PD3ProgressResponse struct {
	Prerequisite PD3PrerequisiteResponse `json:"prerequisite"`
	Objective    PD3ObjectiveResponse    `json:"objective"`
}

type PD3PagingResponse struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
	First    string `json:"first"`
	Last     string `json:"last"`
}
