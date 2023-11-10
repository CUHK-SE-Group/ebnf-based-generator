package graph

// Set is a generic set data structure that uses a map for storage.
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates a new Set from a slice of elements.
func NewSet[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{elements: make(map[T]struct{})}
	for _, e := range elements {
		s.elements[e] = struct{}{}
	}
	return s
}

// Add adds elements to the Set.
func (s *Set[T]) Add(elements ...T) {
	for _, e := range elements {
		s.elements[e] = struct{}{}
	}
}

// Remove removes elements from the Set.
func (s *Set[T]) Remove(elements ...T) {
	for _, e := range elements {
		delete(s.elements, e)
	}
}

// Contains checks if an element is in the Set.
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size returns the number of elements in the Set.
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Elements returns all the elements in the Set as a slice.
func (s *Set[T]) Elements() []T {
	elements := make([]T, 0, len(s.elements))
	for k := range s.elements {
		elements = append(elements, k)
	}
	return elements
}

// Difference returns a new set which is the difference of the two sets (elements in a but not in b).
func Difference[T comparable](a, b *Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for element := range a.elements {
		if _, exists := b.elements[element]; !exists {
			differenceSet.Add(element)
		}
	}
	return differenceSet
}
