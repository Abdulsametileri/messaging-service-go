package viewmodels

type MessageVm struct {
	Text string `json:"text" binding:"required"`
}
