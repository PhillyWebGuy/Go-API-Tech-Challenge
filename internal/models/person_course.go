package models

// PersonCourse represents the many-to-many relationship between Person and Course.
type PersonCourse struct {
	PersonID uint `gorm:"column:person_id;not null" json:"person_id"`
	CourseID uint `gorm:"column:course_id;not null" json:"course_id"`
}

// TableName sets the table name for the PersonCourse struct.
func (PersonCourse) TableName() string {
	return "person_course"
}
