package geo

import (
	"errors"
)

// Point2D defines a node of a slice as an array of 2 coordinates: x and y.
type Point2D [2]float32

// X returns the x coordinate.
func (n Point2D) X() float32 {
	return n[0]
}

// Y returns the y coordinate.
func (n Point2D) Y() float32 {
	return n[1]
}

// Slice defines the resource object for slices.
type Slice struct {
	Vertices []Point2D
	Polygons [][]int
	TopZ     float32
}

// BeginPolygon adds a new polygon and return its index.
func (s *Slice) BeginPolygon() int {
	s.Polygons = append(s.Polygons, make([]int, 0))
	return len(s.Polygons) - 1
}

// AddVertex adds a new vertex to the slice and returns its index.
func (s *Slice) AddVertex(x, y float32) int {
	s.Vertices = append(s.Vertices, Point2D{x, y})
	return len(s.Vertices) - 1
}

// AddPolygonIndex adds a new index to the polygon.
func (s *Slice) AddPolygonIndex(polygonIndex, index int) error {
	if polygonIndex >= len(s.Polygons) {
		return errors.New("invalid polygon index")
	}

	if index >= len(s.Vertices) {
		return errors.New("invalid slice segment index")
	}

	p := s.Polygons[polygonIndex]
	if len(p) > 0 && p[len(p)-1] == index {
		return errors.New("duplicated slice segment index")
	}
	s.Polygons[polygonIndex] = append(s.Polygons[polygonIndex], index)
	return nil
}

// AllPolygonsAreClosed returns true if all the polygons are closed.
func (s *Slice) AllPolygonsAreClosed() bool {
	for _, p := range s.Polygons {
		if len(p) > 1 && p[0] != p[len(p)-1] {
			return false
		}
	}
	return true
}

// IsPolygonValid returns true if the polygon is valid.
func (s *Slice) IsPolygonValid(index int) bool {
	if index >= len(s.Polygons) {
		return false
	}
	p := s.Polygons[index]
	return len(p) > 2
}
