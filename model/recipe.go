package model
import "time"

type Recipe struct {
	RecipeId int					`json:"recipeId"`
	Name string 					`json:"name"`
	Description string				`json:"description"`
	TargetTempCelcius float64		`json:"targetTempCelcius"`
	DurationMinutes time.Duration	`json:"durationMinutes"`
}

type RecipeRun struct {
	RunId int						`json:"runId"`
	StartTime time.Time				`json:"startTime"`
	EndTime time.Time				`json:"endTime"`
	Comments string					`json:"comments"`
	Recipe							`json:"recipe"`
}