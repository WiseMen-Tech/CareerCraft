package models

type Job struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	SkillsReq    []string `json:"skills_required"`
	EducationReq string   `json:"education_req"`
	Location     string   `json:"location"`
	DeadlineDays int      `json:"deadline_days"`
	Sector       string   `json:"sector"`
}
