package game

import "strconv"

/**
 * The Snake game is space oriented, in Cartesian coordinates.
 *
 * Here are defined the positional and movements aspects of the game in Cartesian
 * space.
 */

/**
 * A Vector is the base Cartesian representation used for location of all aspects
 * of the game, extended by other positional elements.
 */

// Vector in space
type Vector struct {
	X, Y int
}

// Detect the same point
func (v Vector) Equals(w Vector) bool {
	return v.X == w.X && v.Y == w.Y
}

// Convert to a printable string
func (v Vector) String() string {
	return "(" + strconv.Itoa(v.X) + "," + strconv.Itoa(v.Y) + ")"
}

/**
 * Unit vectors are directional incrementers which are used to represent viable
 * unit Cartesian movements in the grid.
 */

// Some unit vectors
var (
	Up    Vector = Vector{X: 0, Y: 1}
	Right Vector = Vector{X: 1, Y: 0}
	Down  Vector = Vector{X: 0, Y: -1}
	Left  Vector = Vector{X: -1, Y: 0}
)

/**
 * A Grid is a Cartesian playspace, with unsigned integer dimensions
 * which is represented as a single Vector
 */

// Grid space for snake game
// @NOTE we don't enforce positive dimension on the grid, the game does that.
type Grid Vector

// Convert to a printable string
func (g Grid) String() string {
	return Vector(g).String()
}

// Is a point inside the grid
func (g Grid) Contains(p Point) bool {
	return p.X >= 0 &&
		p.Y >= 0 &&
		p.X <= g.X &&
		p.Y <= g.Y
}

/**
 * A Point is a positional Vector for a point on a grid.  It has no grid awareness
 */

// Point in space struct
type Point Vector

// Detect the same point
func (p Point) Equals(q Point) bool {
	return Vector(p).Equals(Vector(q))
}

// Convert to a printable string
func (p Point) String() string {
	return Vector(p).String()
}

// Move to a new Point from a source Point and a Vector
func (p *Point) Move(v Vector) Point {
	return Point{X: p.X + v.X, Y: p.Y + v.Y}
}
