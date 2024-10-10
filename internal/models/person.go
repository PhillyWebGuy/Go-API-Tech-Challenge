package models

// Person represents a person with associated courses.
type Person struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"column:first_name;not null" validate:"required" json:"first_name"`
	LastName  string `gorm:"column:last_name;not null" validate:"required" json:"last_name"`
	Type      string `gorm:"column:type;not null;check:type IN ('professor', 'student')" validate:"required" json:"type"`
	Age       int    `gorm:"column:age;not null" validate:"required" json:"age"`
}

type PersonWithCourses struct {
	Person
	Courses []int `gorm:"many2many:person_course" json:"courses"`
}

// TableName sets the table name for the Person struct.
func (Person) TableName() string {
	return "person"
}
