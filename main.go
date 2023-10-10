package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rgdotplot/simple-workout-log-to-hevy/data"
	_ "modernc.org/sqlite"
	"os"
	"strconv"
	"strings"
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

	fmt.Println("Done!")
}

func process(input, output string) (err error) {
	_ = os.Remove(output)

	workouts, err := importAll(input)
	if err != nil {
		return
	}

	grouped := groupByDate(workouts)

	var export []data.StrongWorkout
	for _, group := range grouped {
		var groups []data.StrongWorkout

		groups, err = processGroup(group)
		if err != nil {
			return
		}

		export = append(export, groups...)
	}

	// Open a file for writing
	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	defer file.Close()

	// Write heading line
	_, err = file.WriteString("Date;Workout Name;Exercise Name;Set Order;Weight;Weight Unit;Reps;RPE;Distance;Distance Unit;Seconds;Notes;Workout Notes;Workout Duration\n")
	for _, row := range export {
		_, err = file.WriteString(fmt.Sprintf("%s;\"%s\";\"%s\";%d;%.2f;%s;%d;%s;%d;%s;%s;%s;%s;%s\n",
			row.Date,
			row.WorkoutName,
			row.ExerciseName,
			row.SetOrder,
			row.Weight,
			row.WeightUnit,
			row.Reps,
			row.RPE,
			row.Distance,
			row.DistanceUnit,
			row.Seconds,
			row.Notes,
			row.WorkoutNotes,
			row.WorkoutDuration,
		))
		if err != nil {
			return
		}
	}

	return nil
}

func processGroup(workouts []data.SwlWorkout) (list []data.StrongWorkout, err error) {
	var first = workouts[0]
	var last = workouts[len(workouts)-1]

	firstRawTime, err := timeToRaw(first.Time)
	if err != nil {
		return
	}

	lastRawTime, err := timeToRaw(last.Time)
	if err != nil {
		return
	}

	diff := lastRawTime - firstRawTime
	for _, workout := range workouts {
		for i, rep := range workout.Reps {

			list = append(list, data.StrongWorkout{
				Date:            fmt.Sprintf("%s %s:00", first.Date, first.Time),
				WorkoutName:     fmt.Sprintf("Workout %s %s", first.Date, first.Time),
				ExerciseName:    workout.Exercise,
				SetOrder:        i + 1,
				Weight:          rep.Weight,
				WeightUnit:      "kg",
				Reps:            rep.Rep,
				RPE:             "",
				Distance:        0,
				DistanceUnit:    "",
				Seconds:         "",
				Notes:           "",
				WorkoutNotes:    workout.Comment,
				WorkoutDuration: fmt.Sprintf("%ds", diff*60),
			})
		}
	}

	return
}

func timeToRaw(time string) (_ int, err error) {
	timeSep := strings.Split(time, ":")

	hour, err := strconv.Atoi(timeSep[0])
	if err != nil {
		return
	}

	minute, err := strconv.Atoi(timeSep[1])
	if err != nil {
		return
	}

	return hour*60 + minute, nil
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
