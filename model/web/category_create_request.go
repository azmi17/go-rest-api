package web

// pattern method: (Model|Action|Req)
type CategoryCreateRequest struct {
	Name string `validate:"required,min=4,max=200" json:"name"` // tidak butuh (id Int) karena set auto-increment
}
