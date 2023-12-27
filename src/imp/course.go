package imp

type Course struct {
	CourseName string
	TimeSlot   string
}

func NewCourse(courseName string, timeSlot string) *Course {
	return &Course{
		CourseName: courseName,
		TimeSlot:   timeSlot,
	}
}

func (c *Course) DeepCopy() Course {
	return *NewCourse(c.CourseName, c.TimeSlot)
}

func (c *Course) Equals(other *Course) bool {
	return c.CourseName == other.CourseName && c.TimeSlot == other.TimeSlot
}
