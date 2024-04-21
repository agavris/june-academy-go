package imp

import "fmt"

type Section struct {
	Course      *Course
	MaxStudents int
	Students    []*Student
}

var events = map[string]int{
	"Our Feathered Friends":                                      24,
	"Songwriters' Seminar":                                       24,
	"Quilting for a Cause":                                       10,
	"Historical Massachusetts Sites":                             24,
	"Wellness Walks":                                             24,
	"Ultimate Frisbee":                                           24,
	"Dungeons and Dragons":                                       18,
	"Crash Course: College Essay":                                24,
	"Art History at the MFA - Full Day Course":                   24,
	"Intro to Italian":                                           24,
	"Gods and Heroes: Unveiling the Myths of Rome and Greece":    24,
	"Fabulous Fungi!!!":                                          24,
	"Amazing Race - Full Day Course - JUNIORS ONLY":              48,
	"Web Development and Design":                                 24,
	"Beautiful Morning Walks":                                    24,
	"Crochet and Knit":                                           24,
	"Teamwork and collaboration: Competitive Multiplayer Gaming": 24,
	"Local Hidden History: Public History and the Golden Ball Tavern Museum":         12,
	"Performing Arts Trip (NYC) - Full Day Course":                                   75,
	"Beyond Newbury Street: Neighborhoods of Boston - Full Day Course- JUNIORS ONLY": 18,
	"College Tours - Full Day Course":                                                45,
	"Feminist Beautification and Jewelry Making":                                     24,
	"ABCs of Word Puzzles - Wordle and Beyond":                                       24,
	"Weston's Ecology - Exploring Local Conservation Land":                           24,
	"\"Not Bored, Just Board\" Games":                                                24,
	"Magic The Gathering: Learn to Play and Compete":                                 24,
	"Brain and Body Exercise":                                                        24,
	"Play Like a Kid":                                                                24,
	"Disc Golf":                                                                      24,
	"\"Old School\" Photography":                                                     16,
	"Chemistry of Cooking":                                                           15,
	"Model UN Global Leadership Course on Promoting Democracy Worldwide":             24,
	"Fishing at the Pond":                                                            25,
	"Farming to Benefit the Community":                                               15,
	"Amigurumi: Crocheting Small Stuffed Animals":                                    25,
	"Bridge: Bill Gates's Favorite Card Game 6":                                      25,
	"How to Speak K-Pop":                                                             18,
	"Moth-Inspired Storytelling":                                                     25,
	"Financial Literacy":                                                             25,
	"Boston Impact-Self Defense for Post-High School Life":                           18,
	"Building a Personal Fitness Plan":                                               25,
	"Introduction to Podcasting":                                                     25,
}

func NewSection(course *Course, maxStudents int) *Section {
	sec := &Section{
		Course:      course,
		MaxStudents: maxStudents,
		Students:    make([]*Student, 0),
	}

	if max, ok := events[course.CourseName]; ok {
		sec.MaxStudents = max
	} else {
		// Handle the case where the course name is not found in the map
		fmt.Println("Course name not found in events map. Assigning default value.")
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
