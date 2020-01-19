package game

import (
	"errors"
)

/**
 *  A Game tracks the various coordinated components of a game of Snake.  It
 *  Maintains a Grid space for a playing surface, a Snake, and a Food Point.
 *  The Game is RESPONSIBLE FOR MOVE VALIDATION AND EXECUTION.
 *  That the game needs to implement the following functionality:
 *  1. Validate and run a game Tick, which is clock iteration for moving the
 *     snake forward on its vector.
 *  2. Check if a game Tick produces a Boundary collision or a Snake self Collision
 *  3. Check if a Tick moves a snake onto food, and there grows the snake and
 *     Requires new food.
 *
 *
 *  @NOTE We define NO-FOOD as a FOOD point outside of the Grid, in order to
 *    maintain point memory without a reference. An empty point struct is a valid
 *    in grid food Point.
 */

// Game constructor with more automatic setup
func AutoGame(size Vector, food Point) (Game, error) {
	gr := Grid(size)
	h := Point{X: gr.X / 2, Y: gr.Y / 2}
	s := NewSnake(h, Up) // Default to facing up
	g, err := NewGame(gr, s, food)
	return g, err
}

// NewGame validating Game constructor
func NewGame(gr Grid, s Snake, f Point) (Game, error) {
	g := Game{grid: gr, snake: s, food: f}
	return g, g.Validate()
}

// Game object which can manage a grid and a snake
type Game struct {
	grid  Grid
	snake Snake
	food  Point
}

// Validate the game
func (g *Game) Validate() error {
	// Validate
	if !g.grid.Contains(Point{X: 0, Y: 0}) { // grid is X>0 and Y>0
		return errors.New("Could not create game, grid isn't `positive`.")
	}
	if !g.grid.Contains(g.snake.HeadPoint()) {
		return errors.New("Could not create game, as Snake head point is outside of grid.")
	}
	if !g.grid.Contains(g.food) {
		return errors.New("Could not create game, snake start is outside of the grid.")
	}
	return nil
}

// Get the Grid size as a Vector
func (g *Game) Size() Vector {
	return Vector(g.grid)
}

// Turn to a new direction (Does not step)
// @TODO should we detect turning to the same direction?
func (g *Game) Turn(d Vector) {
	g.snake.Turn(d)
}

// Facing snake direction
func (g *Game) Facing() Vector {
	return g.snake.Facing()
}

// Get the Head segment (can be used for recursion)
func (g *Game) Head() *Segment {
	return g.snake.Head()
}

// Get the HeadPoint
func (g *Game) HeadPoint() Point {
	return g.snake.HeadPoint()
}
func (g *Game) Length() uint {
	return g.snake.Length()
}

// Set a Food Point
func (g *Game) SetFood(f Point) {
	g.food = f
}

// Unset food, like if it was just eaten
func (g *Game) unsetFood() {
	g.food = Point{X: g.grid.X + 1, Y: g.grid.Y + 1}
}

// Get the current Food Point
func (g *Game) Food() (Point, error) {
	if g.NeedsFood() {
		return g.food, errors.New("Game currently has no food.")
	}
	return g.food, nil
}

// Does the game need a new food point
func (g *Game) NeedsFood() bool {
	return !g.grid.Contains(g.food)
}

// Tick the game forward as a step
func (g *Game) Tick() (TickResult, error) {
	shp := g.snake.HeadPoint()
	sd := g.snake.Facing()

	np := shp.Move(sd)

	// Prevent a Move if that will cause a collision, and return an error

	if !g.grid.Contains(np) {
		return TickResult{BoundaryCollision: true}, errors.New("Grid collision")
	}
	if g.snake.Contains(np) {
		return TickResult{SnakeCollision: true}, errors.New("Snake Collision")
	}

	if np.Equals(g.food) {
		g.snake.Grow()
		g.unsetFood()
		return TickResult{AteFood: true, Grew: true}, nil
	} else {
		g.snake.Advance()
		return TickResult{Moved: true}, nil
	}
}

// we could return the results of a step like this
type TickResult struct {
	AteFood           bool // Did the snake eat food (to signal that we need new food)
	Grew              bool // Did the snake grow forward (to signal snake growtch)
	Moved             bool // did the snake move forward
	BoundaryCollision bool // Did the snake collide with the boundary
	SnakeCollision    bool // Did the snake collide with itself (cycle)
}
