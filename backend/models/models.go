package models

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TaskInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Public  bool   `json:"public"`
}

type UpdateTaskInput struct {
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	Completed *bool   `json:"completed"`
	Public    *bool   `json:"public"`
}