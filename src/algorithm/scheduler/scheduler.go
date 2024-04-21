package scheduler

import (
	"encoding/csv"
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm/utils"
	"github.com/agavris/june-academy-go/src/imp"
	"github.com/schollz/progressbar/v3"
	"math/rand"
	"os"
	"strings"
	"time"
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
	if section == nil {
		panic("Attempted to add student to a nil section")
	}

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
		priority := student.StudentPriority
		studentsByPriority[priority] = append(studentsByPriority[priority], student)
	}

	for _, group := range studentsByPriority {
		rand.Shuffle(len(group), func(i, j int) {
			group[i], group[j] = group[j], group[i]
		})
	}

	var shuffledStudents []*imp.Student
	for priority := 1; priority <= 3; priority++ {
		shuffledStudents = append(shuffledStudents, studentsByPriority[priority]...)
	}

	s.DataLoader.Students = shuffledStudents
}

func (s *Scheduler) ScoreSchedule() float64 {
	score := 0.0
	studentCopy := make([]*imp.Student, len(s.DataLoader.Students))
	for i, student := range s.DataLoader.Students {
		score += student.SatisfactionScore()
		if s.BestSchedule != nil && score >= s.BestSchedule.Score {
			return score
		}

		studentCopy[i] = student.DeepCopy()
	}

	//if the current BestSchedule's score is lower than the current score, then we have a new best schedule
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

func (s *Scheduler) CourseNameToSectionToSlice() []*imp.Section {
	sections := make([]*imp.Section, 0, len(s.CourseNameToSection))
	for _, section := range s.CourseNameToSection {
		sections = append(sections, section.DeepCopy())
	}
	return sections
}

func (s *Scheduler) Run(numIterations int) *Schedule {
	bar := progressbar.Default(int64(numIterations))
	sectionFolderPath := "sections/"
	resultsFolderPath := "results/"

	if _, err := os.Stat(sectionFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(sectionFolderPath, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	if _, err := os.Stat(resultsFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(resultsFolderPath, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	currentTime := time.Now()

	file, err := os.Create(fmt.Sprintf("%sresults_%s.csv", resultsFolderPath, currentTime.Format("2006-01-02_15-04-05")))
	if err != nil {
		fmt.Println("Error creating file")
		return nil
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	sectionFile, err := os.Create(fmt.Sprintf("%ssections_%s.csv", sectionFolderPath, currentTime.Format("2006-01-02_15-04-05")))
	if err != nil {
		fmt.Println("Error creating section file:", err)
		return nil
	}
	defer sectionFile.Close()

	sectionWriter := csv.NewWriter(sectionFile)
	defer sectionWriter.Flush()

	for range make([]int, numIterations) {
		err := bar.Add(1)
		if err != nil {
			return nil
		}
		s.ExtractByGradeAndShuffle()
		s.AssignStudentsToSections()
		s.ScoreSchedule()
		s.ClearSections()
	}
	fmt.Println(s.BestSchedule.Score)

	err = writer.Write([]string{"Email", "First Name", "Last Name", "Grade", "AM Course", "PM Course", "FD Course", "SS Score"})
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
		return nil
	}

	if err := sectionWriter.Write([]string{"Course Name", "Max Students", "Enrolled Students", "Student Roster"}); err != nil {
		fmt.Println("Error writing section header to CSV file:", err)
		return nil
	}

	for _, student := range s.BestSchedule.Students {
		err := writer.Write([]string{
			student.StudentEmail,
			student.StudentFirstName,
			student.StudentLastName,
			student.Grade,
			student.EnrolledCourses.AMCourse.CourseName,
			student.EnrolledCourses.PMCourse.CourseName,
			student.EnrolledCourses.FullDayCourse.CourseName,
			fmt.Sprintf("%.6f", student.SatisfactionScore()),
		})
		if err != nil {
			fmt.Println("Error writing to CSV file:", err)
			return nil
		}
		fmt.Println(fmt.Sprintf("%s %s, %s, AM: %s, PM: %s, FD: %s, SS: %f", student.StudentFirstName, student.StudentLastName, student.Grade, student.EnrolledCourses.AMCourse.CourseName, student.EnrolledCourses.PMCourse.CourseName, student.EnrolledCourses.FullDayCourse.CourseName, student.SatisfactionScore()))
	}

	for _, section := range s.BestSchedule.Sections {
		var studentNames []string
		for _, student := range section.Students {
			studentNames = append(studentNames, student.StudentFirstName+" "+student.StudentLastName)
		}

		// Join student names with a newline character for better readability in the CSV
		studentRoster := strings.Join(studentNames, "\n")

		if err := sectionWriter.Write([]string{
			section.Course.CourseName,
			fmt.Sprintf("%d", section.MaxStudents),
			fmt.Sprintf("%d", len(section.Students)),
			fmt.Sprintf("\"%s\"", studentRoster), // Ensure names are enclosed in quotes
		}); err != nil {
			fmt.Println("Error writing section data to CSV file:", err)
			return nil
		}
	}

	return s.BestSchedule
}
