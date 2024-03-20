package proto

type StuScores struct {
	StuId  int `json:"stu_id"`
	Score1 int `json:"score_1"`
	Score2 int `json:"score_2"`
	Score3 int `json:"score_3"`
}

type StuAvgScore struct {
	StuId    int     `json:"stu_id"`
	AvgScore float64 `json:"avg_score"`
}
