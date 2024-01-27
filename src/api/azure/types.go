package azure

type (
	TTSRequest struct {
		Lang         string `json:"lang"`
		Name         string `json:"name"`
		Text         string `json:"text"`
		OutputFormat string `json:"outputFormat"`
	}

	TTSResult struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    struct {
			Data string `json:"data"`
		} `json:"data"`
		Message string `json:"message"`
		TraceId string `json:"traceId"`
	}
)
