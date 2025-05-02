package forms

type errors map[string][]string

// Adds a new error message for a give field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Returns the first error message for a given field.
func (e errors) Get(field string) string {
	if len(e[field]) == 0 {
		return ""
	}
	return e[field][0]
}
