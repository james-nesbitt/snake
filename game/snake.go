package game

import "errors"

/**
 * Aspects of the snake as a linked list of points are represented Here
 *
 * A Snake is composed of Segments, which are linked list nodes of points
 * The behaviour is pretty standard Linked List but focused on functionality
 * related to what we need for the snake game.
 *
 * These are:
 *  1. A snake is a combination of a head and a direction
 *  2. Segments of the snake are nodes of a linked List
 *  3. The snake only moves forward, so the Head of the snake is the most important Segment
 *  4. When a snake moves, we set a new Head segment and cut off the last segment
 * 	5. A snake doesn't detect it's own Cycles (the Game does all of the validation)
 *     The snake does have capability to discover if it contains a Point
 *  6. A snake either moves forward, or grows forward
 *  7. A grow means to add a new head which increases length
 *  8. A move means to add a new head, but to drop the last segment
 */

// Segment as a linked list element of a snake
type Segment struct {
	next  *Segment
	point Point
}

// Point for the segment indicating its position
func (s *Segment) Point() Point {
	return s.point
}

// Next segment in the linked list
func (s *Segment) Next() *Segment {
	return s.next
}

// Length of the linked list
func (s *Segment) Length() uint {
	if s.Next() == nil {
		return 1
	} else {
		return 1 + s.Next().Length()
	}
}

// Remove the last segment in the list by unreferencing the second last .next
func (s *Segment) Pop() error {
	if s.Next() == nil {
		return errors.New("List too short to cut off end")
	}

	if s.Next().Next() == nil {
		s.next = nil
	} else {
		s.Next().Pop()
	}

	return nil
}

// Detect Point in Linked List
// @note We never need a full cycle test as we only ever need to test the head
//       Point, as it is the only new point in the snake
func (s *Segment) FindPoint(p Point) bool {
	if s.point.Equals(p) {
		return true
	} else if s.next == nil {
		return false
	} else {
		return s.next.FindPoint(p)
	}
}

// String convert the Segment for logging
func (s *Segment) String() string {
	if s.Next() == nil {
		return s.Point().String()
	} else {
		return s.Point().String() + "," + s.Next().String()
	}
}

// New Snake from a head point and a Vector
func NewSnake(p Point, d Vector) Snake {
	return Snake{head: &Segment{point: p}, dir: d}
}

// Snake as a linked list of segments and a facing direction
type Snake struct {
	head *Segment
	dir  Vector
}

// Get the Head Segment for the snake
func (s *Snake) Head() *Segment {
	return s.head
}

// Get the Head Point for the snake
func (s *Snake) HeadPoint() Point {
	return s.Head().Point()
}

// Get the snake facing direction
func (s *Snake) Facing() Vector {
	return s.dir
}

// How long is the Snake
// @NOTE mainly used for testing
func (s *Snake) Length() uint {
	return s.Head().Length()
}

// Change Direction of the snake, to any vector
// @NOTE we don't confirm that d is a unit vector, meaning we allow any vector
//       for a direction
func (s *Snake) Turn(d Vector) {
	s.dir = d
}

// Grow the snake ahead one step in its direction by adding a new head segment
func (s *Snake) Grow() {
	h := s.Head()
	hp := h.Point()

	nhp := hp.Move(s.dir)
	nh := Segment{next: h, point: nhp}

	s.head = &nh
}

// Move the snake ahead one step in its direction by adding a new head segment
// and removing the last element in the list
func (s *Snake) Advance() {
	s.Grow()
	s.Head().Pop()
}

// Detect if a Point is in the Snake
func (s *Snake) Contains(p Point) bool {
	return s.Head().FindPoint(p)
}

// Points of the snake as a slice
// @NOTE I am not convinced that we should use this as opposed to relying on the
//       the snake head point with recursive functionality
func (s *Snake) Points() []Point {
	ps := []Point{}

	current := s.Head()
	for current != nil {
		ps = append(ps, current.Point())
		current = current.Next()
	}

	return ps
}
