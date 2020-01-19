package game_test

import (
	"github.com/james-nesbitt/snake/game"
	"testing"
)

/**
 * It would be best to avoid vector components as much as possible in this suite
 */

var (
	gridHeight, gridWidth = 10, 10
	headX, headY          = 4, 6
	head                  = game.Point{X: headX, Y: headY}
)

func starter_snake(dir game.Vector) game.Snake {
	return game.NewSnake(head, dir)
}

// Some basic snake testing
func Test_snake_basic(t *testing.T) {
	s := starter_snake(game.Up)

	shp := s.HeadPoint()
	if !shp.Equals(head) {
		t.Errorf("New snake had wrong starter point")
	}
	if !s.Contains(head) {
		t.Errorf("Snake didn't detect that it contains its own head point")
	}

	// point should have no next.
	if s.Head().Next() != nil {
		t.Errorf("New snake thinks it has more than one segment")
	}

	if s.Length() != 1 {
		t.Errorf("New snake reports incorrect length")
	}

	if s.Facing() != game.Up {
		t.Errorf("Snake wasn't facing expected direction: %s", s.Facing())
	}
}

// Test some basic snake movement
func Test_snake_unitmove(t *testing.T) {
	s := starter_snake(game.Up)
	shp := s.HeadPoint()
	sl := s.Length()

	s.Advance()

	sahp := s.HeadPoint()
	tp := shp.Move(game.Up) // determine what the new head point should be

	if !sahp.Equals(tp) {
		t.Errorf("Snake advance landed on the wrong point")
	}
	if sl != s.Length() {
		t.Errorf("Moving snake has incorrect length. Wanted %d, got %d", sl, s.Length())
	}
}

// Test some snake growing
func Test_snake_grow(t *testing.T) {
	s := starter_snake(game.Up)
	shp := s.HeadPoint()
	sl := s.Length()

	s.Grow()

	sahp := s.HeadPoint()
	tp := shp.Move(game.Up) // determine what the new head point should be

	if !sahp.Equals(tp) {
		t.Errorf("Snake advance landed on the wrong point")
	}
	if sl+1 != s.Length() {
		t.Errorf("Growing snake has incorrect length. Wanted %d, got %d", sl+1, s.Length())
	}
}

// Test snake iteration (and Turn/Grow)
func Test_snake_points(t *testing.T) {
	moves := []game.Vector{
		game.Up,
		game.Right,
		game.Down,
	}

	// Create a snake and grow across some expected path

	s := starter_snake(game.Up)
	shp := s.HeadPoint()
	tps := []game.Point{shp} // track a slice of points to match agains
	for _, dir := range moves {
		t.Logf("Growing %s", dir)
		shp = shp.Move(dir)
		tps = append([]game.Point{shp}, tps...) // prepend

		s.Turn(dir)
		s.Grow()
	}

	// Iterate across the expected points and confirm that the snake has only them

	sh := s.Head()
	for i, tp := range tps {
		t.Logf("Testing point %d : %s", i+1, tp)
		if sh == nil {
			t.Errorf("Snake doesn't have enought points")
			break
		}
		sp := sh.Point()

		if !tp.Equals(sp) {
			t.Errorf("Moved snake has incorrect point %d, %s", i, sp)
		}

		sh = sh.Next()
	}
	if sh != nil {
		t.Errorf("Grown snake has too many Points")
	}
}

// Test that turning the snake moves as a Point.Move() (which we have a unit test for)
func Test_snake_turn(t *testing.T) {
	moves := []game.Vector{
		game.Down,
		game.Down,
		game.Left,
		game.Left,
		game.Up,
		game.Up,
		game.Up,
	}

	s := starter_snake(game.Up)

	shp := s.HeadPoint() // get a point to match the snake movement

	// move in some semi-random number of directions
	for _, dir := range moves {
		t.Log("Moving")
		shp = shp.Move(dir)
		s.Turn(dir)
		s.Grow()

		if !s.Contains(shp) {
			t.Errorf("Snake doesn't think it contains a new Point")
		}
	}

	tp := s.HeadPoint()
	if !shp.Equals(tp) {
		t.Errorf("Snake advance landed on the wrong point")
	}
}
