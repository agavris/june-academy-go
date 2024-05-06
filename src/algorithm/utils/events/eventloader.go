package events

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Course struct {
	Name        string
	MaxStudents int
	TimeSlot    string
}

// ReadCourses reads courses from a CSV file and returns a slice of Course
func ReadCourses(filePath string) ([]Course, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, record := range records {
		maxStudents, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Printf("Error converting student count: %v\n", err)
			continue // Skip records with invalid student counts
		}
		courses = append(courses, Course{
			Name:        record[0],
			MaxStudents: maxStudents,
			TimeSlot:    record[2],
		})
	}
	return courses, nil
}

// MapCoursesToMaxStudents maps course names to their maximum number of students
func MapCoursesToMaxStudents(courses []Course) map[string]int {
	courseToMax := make(map[string]int)
	for _, course := range courses {
		courseToMax[course.Name] = course.MaxStudents
	}
	return courseToMax
}

// MapCoursesToTimeSlots maps course names to their time slots
func MapCoursesToTimeSlots(courses []Course) map[string]string {
	courseToTime := make(map[string]string)
	for _, course := range courses {
		courseToTime[course.Name] = course.TimeSlot
	}
	return courseToTime
}
