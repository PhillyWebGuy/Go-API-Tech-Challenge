package models

// Course represents a course with an ID and a name.
type Course struct {
	ID     uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string   `gorm:"column:name;not null" json:"name"`
	People []Person `gorm:"many2many:person_course;" json:"people"`
}

// TableName sets the table name for the Course struct.
func (Course) TableName() string {
	return "course"
}
