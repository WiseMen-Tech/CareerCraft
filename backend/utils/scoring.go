package utils

import (
	"careerconnect/backend/models"
	"strings"
)

// Jaccard similarity for skill matching
func skillMatch(candidateSkills, jobSkills []string) float64 {
	matchCount := 0
	for _, cs := range candidateSkills {
		for _, js := range jobSkills {
			if strings.ToLower(cs) == strings.ToLower(js) {
				matchCount++
			}
		}
	}
	if len(candidateSkills)+len(jobSkills)-matchCount == 0 {
		return 0
	}
	return float64(matchCount) / float64(len(candidateSkills)+len(jobSkills)-matchCount)
}

// Rule-based scoring
func Score(c models.Candidate, j models.Job) float64 {
	skillScore := skillMatch(c.Skills, j.SkillsReq) * 0.45

	eduScore := 0.0
	if strings.ToLower(c.Education) == strings.ToLower(j.EducationReq) {
		eduScore = 1.0 * 0.20
	}

	locScore := 0.0
	if strings.ToLower(c.Location) == strings.ToLower(j.Location) {
		locScore = 1.0 * 0.15
	}

	sectorScore := 0.0
	for _, interest := range c.Interests {
		if strings.ToLower(interest) == strings.ToLower(j.Sector) {
			sectorScore = 1.0 * 0.10
		}
	}

	deadlineScore := 0.0
	if j.DeadlineDays <= 7 {
		deadlineScore = 1.0 * 0.10
	}

	return skillScore + eduScore + locScore + sectorScore + deadlineScore
}
