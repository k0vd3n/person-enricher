package models

// Person — structure returned by GET /people
type Person struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Surname     string `gorm:"type:varchar(100);not null" json:"surname"`
	Patronymic  string `gorm:"type:varchar(100)" json:"patronymic,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `gorm:"type:varchar(10)" json:"gender,omitempty"`
	Nationality string `gorm:"type:varchar(2)" json:"nationality,omitempty"`
}

// CreatePersonRequest — body of POST /people
type CreatePersonRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

// UpdatePersonRequest — body of PUT /people/{id}
type UpdatePersonRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// PeopleFilter — parameters for GET /people?page=&size=
type PeopleFilter struct {
	Filter string
	Page   int
	Size   int
}

// ErrorResponse — single JSON error response
type ErrorResponse struct {
	Error string `json:"error"`
}
