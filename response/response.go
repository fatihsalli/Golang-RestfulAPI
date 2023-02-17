package response

type JSONSuccessResultData struct {
	TotalItemCount int         `json:"totalitemcount"`
	Data           interface{} `json:"data"`
}

type JSONSuccessResultId struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}
