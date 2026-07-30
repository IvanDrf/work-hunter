package models

import "time"

type Currency int

const (
	RUB Currency = 0
	USD Currency = 1
	EUR Currency = 2
)

type RemoteType int

const (
	OFFICE RemoteType = 0
	REMOTE RemoteType = 1
	HYBRID RemoteType = 2
	ANY    RemoteType = 3
)

type TimeType int

type VacancyStatus int

const (
	MODERATING VacancyStatus = 0
	PUBLISHED  VacancyStatus = 1
	CLOSED     VacancyStatus = 2
	DELETED    VacancyStatus = 3
)

const (
	FULL TimeType = 0
	PART TimeType = 1
)

type Vacancy struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Requirements string `json:"requirements,omitempty"`
	Conditions   string `json:"conditions,omitempty"`

	SalaryMin *uint64  `json:"salary_min,omitempty"`
	SalaryMax *uint64  `json:"salary_max,omitempty"`
	Currency  Currency `json:"currency,omitempty"`

	City  *string `json:"city,omitempty"`
	Metro *string `json:"metro,omitempty"`

	RemoteType RemoteType `json:"remote_type,omitempty"`
	TimeType   TimeType   `json:"time_type,omitempty"`

	ExperienceMin *uint32  `json:"experience_min,omitempty"`
	ExperienceMax *uint32  `json:"experience_max,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type VacancyInfo struct {
	Vacancy

	VacancyID uint64 `json:"vacancy_id"`

	IsCityValid  bool `json:"is_city_valid"`
	IsMetroValid bool `json:"is_metro_valid"`

	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`

	Status            VacancyStatus `json:"status"`
	ModeratedAt       *time.Time    `json:"moderated_at,omitempty"`
	ModeratorComments *string       `json:"moderator_comments,omitempty"`

	Views        *uint64 `json:"views,omitempty"`
	Applications *uint64 `json:"applications"`

	AuthorName string `json:"author_name"`
}
