package data

type SwlWorkout struct {
	ID         int    `db:"id"`
	Date       string `db:"date"`
	Time       string `db:"time"`
	ExerciseID int    `db:"exercise_id"`
	Exercise   string `db:"exercise"`
	Type       string `db:"type"`
	Comment    string `db:"comment"`
	Reps       []SwlRep
}

type SwlRep struct {
	ID         int     `db:"id"`
	ExerciseID int     `db:"date_id"`
	Weight     float64 `db:"weight"`
	Rep        int     `db:"rep"`
}
