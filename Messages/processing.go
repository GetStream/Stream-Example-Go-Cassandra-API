package Messages

import (
	"net/http"
)

func processFormField(r *http.Request, field string) (string, string) {
	fieldData := r.PostFormValue(field)
	if len(fieldData) == 0 {
		return "", "Missing '"+field+"' parameter, cannot continue"
	}
	return fieldData, ""
}
