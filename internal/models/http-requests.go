package models

// Person — structure returned by GET /people
type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"` // теперь просто string
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
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
