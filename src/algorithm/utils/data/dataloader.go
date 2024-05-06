package data

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm"
	"github.com/agavris/june-academy-go/src/algorithm/utils/events"
	"github.com/agavris/june-academy-go/src/imp"
	"github.com/gocarina/gocsv"
	"os"
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
		fmt.Println("Also please make sure that you have the correct number of columns and the names are specified correctly \n as shown in the instruction sheet.")
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

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

	for _, request := range d.Requests {
		student := imp.NewStudent(request.FirstName, request.LastName, request.Email, converter[request.Grade], request, request.Grade)
		d.Students = append(d.Students, student)
	}
}

func (d *DataLoader) loadCourses() {
	courseSet := make(map[string]string)
	courses, _ := events.ReadCourses("events.csv")
	coursesToTime := events.MapCoursesToTimeSlots(courses)
	getTimeSlot := func(courseName string, fieldType string) string {
		if time, ok := coursesToTime[courseName]; ok {
			return time
		}
		fmt.Printf("Course name not found in events map. Please check to make sure the names match in both your events.csv file and your jadata.csv file! Course name: %s, Field type: %s\n", courseName, fieldType)
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
