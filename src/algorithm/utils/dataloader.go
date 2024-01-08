package utils

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm"
	"github.com/agavris/june-academy-go/src/imp"
	"github.com/gocarina/gocsv"
	"os"
	"strings"
)

type DataLoader struct {
	Requests []*algorithm.Request
	Students []*imp.Student
	Courses  []*imp.Course
}

func NewDataLoader() *DataLoader {
	loader := &DataLoader{}
	loader.loadData()
	return loader
}

func (d *DataLoader) loadData() {
	d.loadRequests()
	d.loadStudents()
	d.loadCourses()
}

func (d *DataLoader) loadRequests() {
	file, err := os.OpenFile("jadata.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		fmt.Println("Ensure that your file is named jadata.csv and is in the same directory as the executable.")
		return
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, &d.Requests); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		fmt.Println("Ensure that your file is named jadata.csv and is in the same directory as the executable.")
		return
	}
}

func (d *DataLoader) loadStudents() {
	converter := make(map[string]int)
	converter["Freshman"] = 3
	converter["Sophomore"] = 2
	converter["Junior"] = 1

	for i, request := range d.Requests {
		student := imp.NewStudent(i, converter[request.Grade], request, request.Grade)
		d.Students = append(d.Students, student)
	}
}

func (d *DataLoader) loadCourses() {
	courseSet := make(map[string]string)

	getTimeSlot := func(courseName string, fieldType string) string {
		if strings.Contains(courseName, "Full Day") {
			return "FullDay"
		}
		if strings.HasPrefix(fieldType, "AMFD") {
			return "AM"
		}
		if strings.HasPrefix(fieldType, "PM") {
			return "PM"
		}
		return ""
	}
	for _, request := range d.Requests {
		courseFields := map[string]string{
			"AMFD1": request.AMFD1, "AMFD2": request.AMFD2, "AMFD3": request.AMFD3,
			"AMFD4": request.AMFD4, "AMFD5": request.AMFD5,
			"PM1": request.PM1, "PM2": request.PM2, "PM3": request.PM3,
			"PM4": request.PM4, "PM5": request.PM5,
		}
		for fieldType, course := range courseFields {
			if course != "" {
				courseSet[course] = getTimeSlot(course, fieldType)
			}
		}
	}

	for courseName, timeslot := range courseSet {
		course := imp.NewCourse(courseName, timeslot)
		d.Courses = append(d.Courses, course)
	}
}
