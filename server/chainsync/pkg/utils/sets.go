package utils

type StringSet map[string]struct{}

// Add an element to the set
func (s StringSet) Add(element string) {
	s[element] = struct{}{}
}

// Remove an element from the set
func (s StringSet) Remove(element string) {
	delete(s, element)
}

// Check if an element is present in the set
func (s StringSet) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

// Convert the set to a slice of strings
func (s StringSet) Values() []string {
	values := make([]string, 0, len(s))
	for key := range s {
		values = append(values, key)
	}
	return values
}
