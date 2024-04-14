package core

import "time"

type Banner struct {
	Id        int
	Tags      *[]int
	FeatureID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Banner    string
	IsActive  bool
}

type BannerRequest struct {
	TagID           int    `json:"tag_id" validate:"required,numeric,gte=0"`
	FeatureID       int    `json:"feature_id" validate:"required,numeric,gte=0"`
	UseLastRevision bool   `json:"use_last_revision"`
	Token           string `header:"token" validate:"required"`
}

type BannerRequest2 struct {
	TagID     int    `json:"tag_id" validate:"numeric,gte=0"`
	FeatureID int    `json:"feature_id" validate:"numeric,gte=0"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Token     string `header:"token" validate:"required"`
}
