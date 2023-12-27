package scheduler

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm/utils"
	"github.com/agavris/june-academy-go/src/imp"
	"math/rand"
	"sort"
)

type Schedule struct {
	Students []*imp.Student
	Sections []*imp.Section
	Score    float64
}

type Scheduler struct {
	DataLoader          *utils.DataLoader
	CourseNameToSection map[string]*imp.Section
	BestSchedule        *Schedule
}

func NewScheduler() *Scheduler {
	scheduler := &Scheduler{
		DataLoader:          utils.NewDataLoader(),
		CourseNameToSection: make(map[string]*imp.Section),
	}
	scheduler.loadSections()
	return scheduler
}

func (s *Scheduler) loadSections() {
	for _, course := range s.DataLoader.Courses {
		section := imp.NewSection(course, 0)
		s.CourseNameToSection[course.CourseName] = section
	}
}

func (s *Scheduler) safeAddStudentToSection(student *imp.Student, section *imp.Section) bool {

	if len(section.Students) < section.MaxStudents {
		section.AddStudent(student)
		student.AddEnrolledCourse(section.Course)
		return true
	}
	return false
}

func (s *Scheduler) safeRemoveStudentFromSection(student *imp.Student, section *imp.Section) bool {
	if len(section.Students) > 0 {
		section.RemoveStudent(student)
		student.AddEnrolledCourse(section.Course)
		return true
	}
	return false
}

func (s *Scheduler) FindFirstAvailableSectionForStudent(student *imp.Student, timeSlot string) *imp.Section {
	var courseNames []string

	if timeSlot == "AM" || timeSlot == "FullDay" {
		courseNames = student.RequestedCourses.GetAMCourses()
	} else if timeSlot == "PM" {
		courseNames = student.RequestedCourses.GetPMCourses()
	}
	for _, courseName := range courseNames {
		if courseName != "" {
			section := s.CourseNameToSection[courseName]
			if section != nil && len(section.Students) < section.MaxStudents {
				return section
			}
		}
	}
	return s.GetFirstAvailableSectionWithoutRequest(timeSlot)
}

func (s *Scheduler) GetFirstAvailableSectionWithoutRequest(timeSlot string) *imp.Section {
	for _, section := range s.CourseNameToSection {
		if section.Course.TimeSlot == timeSlot && len(section.Students) < section.MaxStudents {
			return section
		}
	}
	return nil
}

func (s *Scheduler) AssignStudentsToSections() {
	assignCourses := func(student *imp.Student, timeSlot string) *imp.Section {
		course := s.FindFirstAvailableSectionForStudent(student, timeSlot)
		if course == nil {
			secondChoice := s.GetFirstAvailableSectionWithoutRequest(timeSlot)
			s.safeAddStudentToSection(student, secondChoice)
			return secondChoice
		}
		s.safeAddStudentToSection(student, course)
		return course
	}
	for _, student := range s.DataLoader.Students {
		am_course := assignCourses(student, "AM")
		if am_course != nil && am_course.Course.TimeSlot == "AM" {
			assignCourses(student, "PM")
		}
	}
}

func (s *Scheduler) ExtractByGradeAndShuffle() {
	studentsByPriority := make(map[int][]*imp.Student, 3)

	for _, student := range s.DataLoader.Students {
		student.UnrollEverything()
		studentsByPriority[student.StudentPriority] = append(studentsByPriority[student.StudentPriority], student)
	}
	keys := make([]int, 0, len(studentsByPriority))
	for key := range studentsByPriority {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, group := range studentsByPriority {
		rand.Shuffle(len(group), func(i, j int) {
			group[i], group[j] = group[j], group[i]
		})
	}

	var shuffledStudents []*imp.Student
	for _, key := range keys {
		shuffledStudents = append(shuffledStudents, studentsByPriority[key]...)
	}

	s.DataLoader.Students = shuffledStudents
}

func (s *Scheduler) ScoreSchedule() float64 {
	score := 0.0
	studentCopy := make([]*imp.Student, len(s.DataLoader.Students))
	for i, student := range s.DataLoader.Students {
		score += student.SatisfactionScore()
		studentCopy[i] = student.DeepCopy()
	}

	//if the current BestSchedule's score is higher than the current score, then we have a new best schedule
	if s.BestSchedule == nil || score < s.BestSchedule.Score {
		s.BestSchedule = &Schedule{
			Students: studentCopy,
			Sections: s.CourseNameToSectionToSlice(),
			Score:    score,
		}
	}
	return score
}

func (s *Scheduler) ClearSections() {
	for _, section := range s.CourseNameToSection {
		section.ClearStudents()
	}
}

func (s *Scheduler) Run(numIterations int) {
	for _ = range make([]int, numIterations) {
		s.ExtractByGradeAndShuffle()
		s.AssignStudentsToSections()
		s.ScoreSchedule()
		s.ClearSections()
	}
	fmt.Println(s.BestSchedule.Score)
	for _, student := range s.BestSchedule.Students {
		fmt.Println(student)
	}
	for _, section := range s.BestSchedule.Sections {
		fmt.Println(section)
	}

}

func (s *Scheduler) PrintEverything() {
	for _, section := range s.CourseNameToSection {
		fmt.Println(section)
		for _, student := range section.Students {
			fmt.Println(student)
		}
	}

	for _, student := range s.DataLoader.Students {
		fmt.Println(student)
	}
}

func (s *Scheduler) CourseNameToSectionToSlice() []*imp.Section {
	var sections []*imp.Section
	for _, section := range s.CourseNameToSection {
		sections = append(sections, section.DeepCopy())
	}
	return sections
}
