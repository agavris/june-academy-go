package imp

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm"
)

type EnrolledCourses struct {
	AMCourse      Course
	PMCourse      Course
	FullDayCourse Course
}

type Student struct {
	StudentID        int
	StudentPriority  int
	Grade            string
	ID               int
	EnrolledCourses  *EnrolledCourses
	RequestedCourses *algorithm.Request
}

func NewStudent(studentID int, studentPriority int, requestedCourses *algorithm.Request, grade string, id int) *Student {
	return &Student{
		StudentID:        studentID,
		StudentPriority:  studentPriority,
		EnrolledCourses:  &EnrolledCourses{},
		RequestedCourses: requestedCourses,
		Grade:            grade,
		ID:               id,
	}
}

func (s *Student) AddEnrolledCourse(course *Course) {
	if course.TimeSlot == "AM" {
		s.EnrolledCourses.AMCourse = *course
	} else if course.TimeSlot == "PM" {
		s.EnrolledCourses.PMCourse = *course
	} else {
		s.EnrolledCourses.FullDayCourse = *course
	}
}

func (s *Student) RemoveEnrolledCourse(course *Course) {
	if course.TimeSlot == "AM" {
		s.EnrolledCourses.AMCourse = Course{}
	} else if course.TimeSlot == "PM" {
		s.EnrolledCourses.PMCourse = Course{}
	} else {
		s.EnrolledCourses.FullDayCourse = Course{}
	}
}

func (s *Student) UnrollEverything() {
	s.EnrolledCourses = &EnrolledCourses{}
}

//func (s *Student) SatisfactionScore() float64 {
//	score := 0.0
//
//	// Function to check if the course is in the requested list and update the score
//	updateScore := func(courseName string, requestedCourses map[string]struct{}, increment float64) {
//		if courseName != "" {
//			if _, found := requestedCourses[courseName]; !found {
//				score += increment
//			}
//		}
//	}
//
//	// Creating maps for AM and PM requested courses
//	requestedAMCourses := make(map[string]struct{})
//	for _, course := range s.RequestedCourses.GetAMCourses() {
//		requestedAMCourses[course] = struct{}{}
//	}
//
//	requestedPMCourses := make(map[string]struct{})
//	for _, course := range s.RequestedCourses.GetPMCourses() {
//		requestedPMCourses[course] = struct{}{}
//	}
//
//	// Check and update score for each course
//	if s.EnrolledCourses.FullDayCourse.CourseName != "" {
//		updateScore(s.EnrolledCourses.FullDayCourse.CourseName, requestedAMCourses, 1)
//		return score
//	}
//
//	updateScore(s.EnrolledCourses.AMCourse.CourseName, requestedAMCourses, 0.5)
//	updateScore(s.EnrolledCourses.PMCourse.CourseName, requestedPMCourses, 0.5)
//
//	return score
//}

//func (s *Student) SatisfactionScore() float64 {
//	score := 0.0
//
//	amCourse := s.EnrolledCourses.AMCourse
//	pmCourse := s.EnrolledCourses.PMCourse
//	fullDayCourse := s.EnrolledCourses.FullDayCourse
//
//	requestedAMCourses := make(map[string]struct{})
//	for _, course := range s.RequestedCourses.GetAMCourses() {
//		requestedAMCourses[course] = struct{}{}
//	}
//
//	requestedPMCourses := make(map[string]struct{})
//	for _, course := range s.RequestedCourses.GetPMCourses() {
//		requestedPMCourses[course] = struct{}{}
//	}
//
//	if fullDayCourse.CourseName != "" {
//		if _, found := requestedAMCourses[fullDayCourse.CourseName]; !found {
//			score += 1
//		}
//		return score
//	}
//
//	if amCourse.CourseName != "" {
//		if _, found := requestedAMCourses[amCourse.CourseName]; !found {
//			score += 0.5
//		}
//	}
//
//	if pmCourse.CourseName != "" {
//		if _, found := requestedPMCourses[pmCourse.CourseName]; !found {
//			score += 0.5
//		}
//	}
//
//	return score
//}

func (s *Student) SatisfactionScore() float64 {
	score := 0.0

	// Function to check if course is in the slice
	isInCourses := func(courseName string, courses []string) bool {
		for _, c := range courses {
			if c == courseName {
				return true
			}
		}
		return false
	}

	// Check FullDayCourse
	fullDayCourse := s.EnrolledCourses.FullDayCourse
	if fullDayCourse.CourseName != "" {
		if !isInCourses(fullDayCourse.CourseName, s.RequestedCourses.GetAMCourses()) {
			score += 1
		}
		return score // Early return if full day course is present
	}

	// Check AMCouse and PMCouse
	amCourse := s.EnrolledCourses.AMCourse
	if amCourse.CourseName != "" && !isInCourses(amCourse.CourseName, s.RequestedCourses.GetAMCourses()) {
		score += 0.5
	}
	pmCourse := s.EnrolledCourses.PMCourse
	if pmCourse.CourseName != "" && !isInCourses(pmCourse.CourseName, s.RequestedCourses.GetPMCourses()) {
		score += 0.5
	}

	return score
}

func (s *Student) CopyEnrolledCourses() *EnrolledCourses {
	if s.EnrolledCourses == nil {
		return &EnrolledCourses{}
	}
	return &EnrolledCourses{
		AMCourse:      s.EnrolledCourses.AMCourse.DeepCopy(),
		PMCourse:      s.EnrolledCourses.PMCourse.DeepCopy(),
		FullDayCourse: s.EnrolledCourses.FullDayCourse.DeepCopy(),
	}
}

func (s *Student) CopyRequestedCourses() *algorithm.Request {
	return &algorithm.Request{
		Grade: s.RequestedCourses.Grade,
		AMFD1: s.RequestedCourses.AMFD1, AMFD2: s.RequestedCourses.AMFD2, AMFD3: s.RequestedCourses.AMFD3, AMFD4: s.RequestedCourses.AMFD4, AMFD5: s.RequestedCourses.AMFD5,
		PM1: s.RequestedCourses.PM1, PM2: s.RequestedCourses.PM2, PM3: s.RequestedCourses.PM3, PM4: s.RequestedCourses.PM4, PM5: s.RequestedCourses.PM5,
	}
}

func (s *Student) DeepCopy() *Student {
	return &Student{
		StudentID:        s.StudentID,
		StudentPriority:  s.StudentPriority,
		EnrolledCourses:  s.CopyEnrolledCourses(),
		RequestedCourses: s.CopyRequestedCourses(),
		Grade:            s.Grade,
		ID:               s.ID,
	}
}

func (s *Student) Equals(other *Student) bool {
	return s.StudentID == other.StudentID
}

func (s *Student) String() string {
	return fmt.Sprintf("ID: %d", s.ID)
	//return fmt.Sprintf("Student ID: %d, Grade: %s, Enrolled AM: %s, Enrolled PM: %s, Enrolled FD: %s, SS: %f", s.StudentID, s.Grade, s.EnrolledCourses.AMCourse.CourseName, s.EnrolledCourses.PMCourse.CourseName, s.EnrolledCourses.FullDayCourse.CourseName, s.SatisfactionScore())
}
