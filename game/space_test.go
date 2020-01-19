package game_test

import (
	"github.com/james-nesbitt/snake/game"
	"testing"
)

var (
	// A testable point
	px, py = 3, 4
	p      = game.Point{X: px, Y: py}
	// a testable vector
	vx, vy = 6, 1
	v      = game.Vector{X: vx, Y: vy}
	// the test point moved by the test vector
	q = p.Move(v)

	g    = game.Grid{Y: 10, X: 10}
	outs = []game.Point{
		{X: -3, Y: 6},
		{X: 5, Y: -2},
		{X: 20, Y: 5},
		{X: 5, Y: 20},
	}
)

// Test that a point behaves like a point
func Test_point_basic(t *testing.T) {
	if p.X != px {
		t.Errorf("Bad Point X, got: %d, want: %d.", p.X, px)
	}
	if p.Y != py {
		t.Errorf("Bad Point Y, got: %d, want: %d.", p.Y, py)
	}
}

// Test point equals
func test_point_equals(t *testing.T) {
	p := game.Point{X: 2, Y: 6}
	pe := game.Point{X: 2, Y: 6}
	pne := []game.Point{
		{X: 1, Y: 1},
		{X: 7, Y: 6},
		{X: 8, Y: 2},
		{X: 2, Y: 5},
		{X: 6, Y: 6},
	}

	if !(p.Equals(pe) && pe.Equals(p)) {
		t.Errorf("Equal points did not think that they are equal")
	}
	for _, ne := range pne {
		if p.Equals(ne) || ne.Equals(p) {
			t.Errorf("Unequal points though that they are equal")
		}
	}
}

// Test that a point'sMove method works
func Test_point_move(t *testing.T) {
	if q.X != px+vx {
		t.Errorf("Bad Point Move X, got: %d, want: %d.", q.X, vx)
	}
	if q.Y != py+vy {
		t.Errorf("Bad Point Move Y, got: %d, want: %d.", q.Y, vy)
	}
}

// Test some unit vector movements
func Test_UnitVectorMoves(t *testing.T) {
	pUp := p.Move(game.Up)
	pLeft := p.Move(game.Left)
	pDown := p.Move(game.Down)
	pRight := p.Move(game.Right)

	if !(pUp.X == p.X && pUp.Y == (p.Y+1)) {
		t.Errorf("Up unit Vector didn't Move up by one")
	}
	if !(pDown.X == p.X && pDown.Y == (p.Y-1)) {
		t.Errorf("Down unit Vector didn't Move down by one")
	}
	if !(pLeft.X == (p.X-1) && pLeft.Y == p.Y) {
		t.Errorf("Left unit Vector didn't Move left by one")
	}
	if !(pRight.X == (p.X+1) && pRight.Y == p.Y) {
		t.Errorf("Right unit Vector didn't Move right by one")
	}
}

// Test the Grid.Contains method
func Test_PointContains(t *testing.T) {
	if !g.Contains(p) {
		t.Errorf("Contains check failed for a point inside the grid")
	}

	for _, out := range outs {
		if g.Contains(out) {
			t.Errorf("Contains check failed to detect that an Point is outside the grid")
		}
	}
}
