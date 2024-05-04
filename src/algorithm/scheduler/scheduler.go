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

//func (s *Scheduler) FindFirstAvailableSectionForStudent(student *imp.Student, timeSlot string) *imp.Section {
//	var courseNames []string
//	priorityCourse := "Model UN Global Leadership Course on Promoting Democracy Worldwide" // Set the priority course name
//
//	if timeSlot == "AM" || timeSlot == "FullDay" {
//		courseNames = student.RequestedCourses.GetAMCourses()
//	} else if timeSlot == "PM" {
//		courseNames = student.RequestedCourses.GetPMCourses()
//	}
//
//	if student.StudentPriority > 1 {
//		for _, courseName := range courseNames {
//			if courseName == priorityCourse {
//				section := s.CourseNameToSection[courseName]
//				if section != nil && len(section.Students) < section.MaxStudents {
//					return section
//				}
//			}
//		}
//	}
//
//	// Then check other courses
//	for _, courseName := range courseNames {
//		if courseName != "" && courseName != priorityCourse {
//			section := s.CourseNameToSection[courseName]
//			if section != nil && len(section.Students) < section.MaxStudents {
//				return section
//			}
//		}
//	}
//
//	// If no requested course is available, fall back to the first available section
//	return s.GetFirstAvailableSectionWithoutRequest(timeSlot)
//}

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
	sectionFolderPath := "sections/"
	resultsFolderPath := "results/"
	currentTime := time.Now()

	// Ensure necessary directories exist
	if err := ensureDirectory(sectionFolderPath); err != nil {
		fmt.Println("Failed to create section directory:", err)
		return nil
	}
	if err := ensureDirectory(resultsFolderPath); err != nil {
		fmt.Println("Failed to create results directory:", err)
		return nil
	}

	// Setup CSV files for results and sections
	resultsFile, err := setupCSVFile(resultsFolderPath, "results_", currentTime)
	if err != nil {
		fmt.Println("Error setting up results CSV file:", err)
		return nil
	}
	defer resultsFile.Close()
	resultsWriter := csv.NewWriter(resultsFile)
	defer resultsWriter.Flush()

	sectionFile, err := setupCSVFile(sectionFolderPath, "sections_", currentTime)
	if err != nil {
		fmt.Println("Error setting up section CSV file:", err)
		return nil
	}
	defer sectionFile.Close()
	sectionWriter := csv.NewWriter(sectionFile)
	defer sectionWriter.Flush()

	// Initialize progress bar for tracking
	bar := progressbar.Default(int64(numIterations))

	for i := 0; i < numIterations; i++ {
		bar.Add(1)
		s.ExtractByGradeAndShuffle()
		s.AssignStudentsToSections()
		s.ScoreSchedule()
		s.ClearSections()
	}

	// Output schedule and section information to CSV files
	if err := outputSchedule(resultsWriter, sectionWriter, s.BestSchedule); err != nil {
		fmt.Println("Error writing to CSV file:", err)
		return nil
	}

	fmt.Println("Best schedule score:", s.BestSchedule.Score)
	return s.BestSchedule
}

func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory doesn't exist, create it
		return os.Mkdir(path, os.ModePerm)
	}
	return nil // Directory already exists
}

func setupCSVFile(path, prefix string, currentTime time.Time) (*os.File, error) {
	filename := fmt.Sprintf("%s%s%s.csv", path, prefix, currentTime.Format("2006-01-02_15-04-05"))
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, create it
		return os.Create(filename)
	}
	// File already exists (this might be ideal if you want to append, but you would need to adjust other code)
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
}

func outputSchedule(resultsWriter, sectionWriter *csv.Writer, schedule *Schedule) error {
	// Writing header rows
	if err := resultsWriter.Write([]string{"Email", "First Name", "Last Name", "Grade", "AM Course", "PM Course", "FD Course", "SS Score"}); err != nil {
		return err
	}
	if err := sectionWriter.Write([]string{"Course Name", "Max Students", "Enrolled Students", "Student Roster"}); err != nil {
		return err
	}

	// Writing data rows
	for _, student := range schedule.Students {
		record := []string{
			student.StudentEmail,
			student.StudentFirstName,
			student.StudentLastName,
			student.Grade,
			student.EnrolledCourses.AMCourse.CourseName,
			student.EnrolledCourses.PMCourse.CourseName,
			student.EnrolledCourses.FullDayCourse.CourseName,
			fmt.Sprintf("%.6f", student.SatisfactionScore()),
		}
		if err := resultsWriter.Write(record); err != nil {
			return err
		}
	}

	for _, section := range schedule.Sections {
		studentNames := make([]string, len(section.Students))
		for i, student := range section.Students {
			studentNames[i] = student.StudentFirstName + " " + student.StudentLastName
		}
		studentRoster := strings.Join(studentNames, ", ")
		record := []string{
			section.Course.CourseName,
			fmt.Sprintf("%d", section.MaxStudents),
			fmt.Sprintf("%d", len(section.Students)),
			fmt.Sprintf("\"%s\"", studentRoster),
		}
		if err := sectionWriter.Write(record); err != nil {
			return err
		}
	}
	return nil
}
