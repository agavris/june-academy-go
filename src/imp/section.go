package imp

import "fmt"

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
	if course.TimeSlot == "AM" || course.TimeSlot == "PM" {
		sec.MaxStudents = 25
	} else {
		sec.MaxStudents = 42
	}
	return sec
}

func (s *Section) AddStudent(student *Student) {
	s.Students = append(s.Students, student)
}

func (s *Section) RemoveStudent(student *Student) {
	for i, st := range s.Students {
		if st.StudentID == student.StudentID {
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
	return fmt.Sprintf("Section: %s, Timeslot: %s, Students: %d", s.Course.CourseName, s.Course.TimeSlot, len(s.Students))
}

func (s *Section) Equals(other *Section) bool {
	return s.Course.Equals(other.Course)
}

func (s *Section) NotEquals(other *Section) bool {
	return !s.Equals(other)
}
