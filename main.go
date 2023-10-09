package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rgdotplot/simple-workout-log-to-hevy/data"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	input := flag.String("input", "input.bak", "input file")
	output := flag.String("output", "output.csv", "output file")

	flag.Parse()

	fmt.Println("input:", *input)
	fmt.Println("output:", *output)

	err := process(*input, *output)
	if err != nil {
		panic(err)
	}
}

func process(input, output string) (err error) {
	_ = os.Remove(output)

	workouts, err := importAll(input)
	if err != nil {
		return
	}

	grouped := groupByDate(workouts)

	fmt.Println(grouped)

	return nil
}

// A function that takes a list of workouts and returns them grouped by date
func groupByDate(workouts []data.SwlWorkout) (grouped map[string][]data.SwlWorkout) {
	grouped = make(map[string][]data.SwlWorkout)

	for _, workout := range workouts {
		grouped[workout.Date] = append(grouped[workout.Date], workout)
	}

	return
}

func importAll(input string) (workouts []data.SwlWorkout, err error) {
	db, err := sqlx.Open("sqlite", input)
	if err != nil {
		return
	}

	defer db.Close()

	err = db.Select(&workouts, "SELECT id, date, time, exercise_id, exercise, type, comment FROM workouts ORDER BY date, time")
	if err != nil {
		return
	}

	for i, workout := range workouts {
		err = db.Select(&workouts[i].Reps, "SELECT id, date_id, weight, rep FROM reps WHERE date_id = ? ORDER BY id", workout.ID)
		if err != nil {
			return
		}
	}

	return
}
