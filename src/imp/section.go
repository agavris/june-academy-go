package imp

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Section struct {
	Course      *Course
	MaxStudents int
	Students    []*Student
}

func loadEventsFromCSV(filePath string) (map[string]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Use bufio.Scanner to handle different types of newline characters
	scanner := bufio.NewScanner(file)
	events := make(map[string]int)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		// Use strings.NewReader to create a reader from the scanned text
		reader := csv.NewReader(strings.NewReader(line))
		reader.Comma = ';'
		reader.LazyQuotes = true // Allow quotes in fields
		reader.TrimLeadingSpace = true

		record, err := reader.Read()
		if err != nil {
			fmt.Printf("Error reading line %d: %v\n", lineCount, err)
			continue // Skip this record and move to the next
		}

		if len(record) != 2 {
			fmt.Printf("Skipping malformed record on line %d (expected 2 fields, got %d): %v\n", lineCount, len(record), record)
			continue
		}

		courseName := strings.Trim(record[0], "\"") // Clean up the course name
		maxStudents, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Printf("Skipping record with invalid number of students on line %d: %v, error: %v\n", lineCount, record, err)
			continue
		}
		events[courseName] = maxStudents
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning file: %w", err)
	}

	return events, nil
}

func NewSection(course *Course, maxStudents int) *Section {
	sec := &Section{
		Course:      course,
		MaxStudents: maxStudents,
		Students:    make([]*Student, 0),
	}

	var err error
	events, err := loadEventsFromCSV("events.csv")
	if err != nil {
		fmt.Println("Error loading events from CSV file: ", err)
	}
	if max, ok := events[course.CourseName]; ok {
		sec.MaxStudents = max
	} else {
		// Handle the case where the course name is not found in the map
		fmt.Println("Course name not found in events map. Please check to make sure the names match in both your events.csv file and your jadata.csv file!")
		panic(course.CourseName)
	}

	return sec
}

func (s *Section) AddStudent(student *Student) {
	s.Students = append(s.Students, student)
}

func (s *Section) RemoveStudent(student *Student) {
	for i, st := range s.Students {
		if st.StudentEmail == student.StudentEmail {
			s.Students = append(s.Students[:i], s.Students[i+1:]...)
			return
		}
	}
}

func (s *Section) ClearStudents() {
	s.Students = make([]*Student, 0)
}

func (s *Section) DeepCopy() *Section {
	studentArray := make([]*Student, len(s.Students))
	for i, student := range s.Students {
		studentArray[i] = student.DeepCopy()
	}

	return &Section{
		Course:      s.Course,
		MaxStudents: s.MaxStudents,
		Students:    studentArray,
	}
}

func (s *Section) String() string {
	//return fmt.Sprintf("Section: %s, Timeslot: %s, Students: %d, List: %s", s.Course.CourseName, s.Course.TimeSlot, len(s.Students), s.Students)
	return fmt.Sprintf("Section: %s, Students: %d", s.Course.CourseName, len(s.Students))
}

func (s *Section) Equals(other *Section) bool {
	return s.Course.Equals(other.Course)
}

func (s *Section) NotEquals(other *Section) bool {
	return !s.Equals(other)
}
