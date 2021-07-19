// error module contains _error corner case processing.
package command

func applyError(v map[string]interface{}, value interface{}) map[string]interface{} {
	vMap := value.(map[string]interface{})
	v["schema"] = map[string]interface{}{
		"properties": map[string]interface{}{
			"error": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code":    map[string]interface{}{"type": "string", "example": vMap["code"]},
					"message": map[string]interface{}{"type": "string", "example": vMap["message"]},
					"status":  map[string]interface{}{"type": "string", "example": vMap["status"]},
				},
			},
		},
	}
	return v
}
