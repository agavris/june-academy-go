package imp

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm/utils/events"
)

type Section struct {
	Course      *Course
	MaxStudents int
	Students    []*Student
}

func NewSection(course *Course, maxStudents int) *Section {
	sec := &Section{
		Course:      course,
		MaxStudents: maxStudents,
		Students:    make([]*Student, 0),
	}

	var err error
	courses, err := events.ReadCourses("events.csv")
	coursesToMax := events.MapCoursesToMaxStudents(courses)
	if err != nil {
		fmt.Println("Error loading events from CSV file: ", err)
	}
	if max, ok := coursesToMax[course.CourseName]; ok {
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
