package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"CareerConnect/models"
	"CareerConnect/utils"
	"sort"
)

// Mock dataset
var Jobs = []models.Job{
	{ID: "1", Title: "Data Entry Intern", SkillsReq: []string{"Excel", "Typing"}, EducationReq: "Undergraduate", Location: "Lucknow", DeadlineDays: 5, Sector: "Administration"},
	{ID: "2", Title: "Web Dev Intern", SkillsReq: []string{"HTML", "CSS", "Go"}, EducationReq: "Undergraduate", Location: "Delhi", DeadlineDays: 10, Sector: "IT"},
	{ID: "3", Title: "AI Research Intern", SkillsReq: []string{"Python", "ML"}, EducationReq: "Postgraduate", Location: "Remote", DeadlineDays: 3, Sector: "Research"},
	{ID: "4", Title: "Govt Clerk Job", SkillsReq: []string{"Typing", "MS Office"}, EducationReq: "Undergraduate", Location: "Lucknow", DeadlineDays: 7, Sector: "Administration"},
}

// Temporary store for candidates
var candidates []models.Candidate

func RecommendHandler(c *gin.Context) {
	var candidate models.Candidate
	if err := c.BindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	candidates = append(candidates, candidate)

	type ScoredJob struct {
		models.Job
		Score float64 `json:"score"`
	}

	var scored []ScoredJob
	for _, job := range Jobs {
		score := utils.Score(candidate, job)
		if score > 0 {
			scored = append(scored, ScoredJob{job, score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	limit := 3
	if len(scored) < 3 {
		limit = len(scored)
	}

	c.JSON(http.StatusOK, scored[:limit])
}

// Check deadlines daily and send alerts
func CheckDeadlines() {
	for _, user := range candidates {
		for _, job := range Jobs {
			if job.DeadlineDays <= 7 {
				utils.SendNotification(user, job, job.DeadlineDays)
			}
		}
	}
}
