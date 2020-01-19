package game_test

import (
	"github.com/james-nesbitt/snake/game"
	"testing"
)

// Test the game constructor (the auto constructor, not the parametrized)
func Test_GameConstruction(t *testing.T) {
	gs := game.Vector{X: 10, Y: 10} // grid size.
	gf := game.Point{X: 3, Y: 4}    // game food point
	g, err := game.AutoGame(gs, gf)

	if err != nil {
		t.Errorf("Error creating game: %s", err)
	}

	if !gs.Equals(g.Size()) {
		t.Errorf("Auto Game does not have the provided size")
	}

	if g.NeedsFood() {
		t.Errorf("Game reports that it needs food when it shouldn't")
	} else if f, err := g.Food(); err != nil {
		t.Errorf("AutoGame did not correctly create food: %s", err)
	} else if !f.Equals(gf) {
		t.Errorf("Auto Game had food in an unexpected location")
	}
}

// Test that a game fails at points
func Test_WonkyGameConstruct(t *testing.T) {
	// null sanity test
	if _, err := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 3, Y: 4}); err != nil {
		t.Errorf("Error creating game: %s", err)
	}

	if _, err := game.AutoGame(game.Vector{X: -1, Y: 10}, game.Point{X: 3, Y: 4}); err == nil {
		t.Errorf("AutoGame constructor did not produce an error on invalid grid (-X)")
	}
	if _, err := game.AutoGame(game.Vector{X: -10, Y: -4}, game.Point{X: 3, Y: 4}); err == nil {
		t.Errorf("AutoGame constructor did not produce an error on invalid grid (-Y)")
	}
	if _, err := game.AutoGame(game.Vector{X: 2, Y: 2}, game.Point{X: 3, Y: 4}); err == nil {
		t.Errorf("AutoGame constructor did not produce an error with food outside of the grid")
	}

}

// Make a testing game object
func testingGame(t *testing.T) TestingGame {
	grid := game.Vector{X: 10, Y: 10}
	food := game.Point{X: 3, Y: 4}
	game, err := game.AutoGame(grid, food)
	if err != nil {
		t.Errorf("Error creating game: %s", err)
	}
	return TestingGame{game: &game, t: t}
}

func Test_GameEat(t *testing.T) {
	gr := game.Vector{X: 10, Y: 10}
	f := game.Point{X: 5, Y: 7}
	g, _ := game.AutoGame(gr, f)
	tg := TestingGame{game: &g, t: t}
	if !g.HeadPoint().Equals(game.Point{X: 5, Y: 5}) {
		// This is not hard defined behaviour but it is needed for this test
		t.Error("Game started with snake in unexpected place")
	} else {
		tg.turn(game.Up) // make sure we are pointing up
		tg.move(1)

		if g.NeedsFood() {
			t.Errorf("Game reprots that it needs food when it should not")
		}

		tg.eat()

		if !g.HeadPoint().Equals(f) {
			t.Error("Snake that ate food does not report its head is on the food point")
		}
		if !g.NeedsFood() {
			t.Errorf("Game should need food after eating, but doesn't report so")
		}
		if nf, err := g.Food(); err == nil {
			t.Error("Game where the snake just ate food still reports food, before new food was created")
		} else if nf.Equals(f) {
			t.Errorf("Game did not remove food after eating it")
		}

		nf := game.Point{X: 2, Y: 3} // New Food Point
		tg.food(nf.X, nf.Y)
		if g.NeedsFood() {
			t.Errorf("Game reports that it needs food right after we set it")
		} else if gf, err := g.Food(); err != nil {
			t.Errorf("Game reported no food after setting new food location")
		} else if gf.Equals(f) {
			t.Error("Game where new food was create does not have food in the right place")
		}
	}
}

// Run a mock game to test a more moves and events
// - a simple game where we move about suceesfully growing to 10 in length
func Test_GamePlaySimple(t *testing.T) {
	tg := testingGame(t)

	tg.move(3)          // (5,5) -> (5,8)
	tg.turn(game.Left)  // turn left
	tg.move(2)          // (5,8) -> 3,8)
	tg.turn(game.Down)  // turn down
	tg.move(3)          // (3,8) -> (3,5)
	tg.eat()            // (3,5) -> (3,4)::(3,5) [2]
	tg.food(7, 8)       // food at (7,8)
	tg.move(2)          // (3,5) -> (3,2)::(3,3)
	tg.turn(game.Right) // turn right
	tg.move(4)          // (3,2) -> (7,2)::(6,2)
	tg.turn(game.Up)    // turn up
	tg.move(5)          // (7,2) -> (7,7)::(7,6)
	tg.eat()            // (7,7) -> (7,8)::(7,6) [3]
	tg.food(1, 8)       // food at (1,8)
	tg.move(1)          // (7,8) -> (7,9)::(7,7)
	tg.turn(game.Left)  // turn left
	tg.move(6)          // (7,9) -> (1,9)::(3,9)
	tg.turn(game.Down)  // turn down
	tg.eat()            // (1,9) -> (1,8)::3,9) [4]
	tg.food(5, 3)       // food at (5,3)
	tg.move(6)          // (1,8) -> 1,2)::(1,5)
	tg.turn(game.Right) // turn right
	tg.move(1)          // (1,2) -> (2,2)::(1,4)
	tg.turn(game.Up)    // turn up
	tg.move(1)          // (2,2) -> (2,3)::(1,3)
	tg.turn(game.Right) // turn right
	tg.move(2)          // (2,3) -> (4,3)::(2,2)
	tg.eat()            // (4,3) -> (5,3)::(2,2) [5]
	tg.food(9, 1)       // food at (9,1)
	tg.move(4)          // (5,3) -> (9,3)::(4,3)
	tg.turn(game.Down)  // turn down
	tg.move(1)          // (9,3) -> (9,2)::(6,3)
	tg.eat()            // (9,2) -> (9,1)::(6,3) [6]
	tg.food(0, 0)       // food at (0,0)
	tg.move(1)          // (9,1) -> (9,0)::(7,3)
	tg.turn(game.Left)  // turn left
	tg.move(8)          // (9,0) -> (1,0)::(6,0)
	tg.eat()            // (1,0) -> (0,0)::(6,0) [7]
	tg.food(6, 7)       // food at (6,7)
	tg.turn(game.Up)    // turn up
	tg.move(7)          // (0,0) -> (0,7)::(0,1)
	tg.turn(game.Right) // turn right
	tg.move(5)          // (0,7) -> (5,7)::(0,6)
	tg.eat()            // (5,7) -> (6,7)::(0,6) [8]
	tg.food(10, 7)      // food at (10,7)
	tg.move(3)          // (6,7) -> (9,7)::(2,7)
	tg.eat()            // (9,7) -> (10,7)::(2,7) [9]
	tg.food(3, 6)       // food at (3,6)
	tg.turn(game.Up)    // turn up
	tg.move(1)          // (10,7) -> (10,8)::(3,7)
	tg.turn(game.Left)  // turn right
	tg.move(7)          // (10,8) -> (3,8)::(10,7)
	tg.turn(game.Down)  // turn down
	tg.move(1)          // (3,8) -> (3,7)::(10,8)
	tg.eat()            // (3,6) -> (3,6)::(9.8) [10]

}

// Test an expected collision without making it too simple
func Test_GamePlayCollision(t *testing.T) {
	tg := testingGame(t)

	tg.move(3)          // (5,5) -> (5,8)
	tg.turn(game.Left)  // turn left
	tg.move(2)          // (5,8) -> 3,8)
	tg.turn(game.Down)  // turn down
	tg.move(3)          // (3,8) -> (3,5)
	tg.eat()            // (3,5) -> (3,4)::(3,5) [2]
	tg.food(3, 3)       // food at (3,3)
	tg.eat()            // (3,3) -> (3,3)::(3,5) [3]
	tg.food(3, 2)       // food at (3,2)
	tg.eat()            // (3,2) -> (3,2)::(3,5) [4]
	tg.food(3, 1)       // food at (3,1)
	tg.eat()            // (3,2) -> (3,1)::(3,5) [5]
	tg.food(3, 0)       // food at (3,0)
	tg.eat()            // (3,1) -> (3,0)::(3,5) [6]
	tg.food(4, 0)       // food at (4,0)
	tg.turn(game.Right) // turn right
	tg.eat()            // (3,0) -> (4,0)::(3,5) [7]
	tg.food(5, 0)       // food at (5,0)
	tg.eat()            // (4,0) -> (5,0)::(3,5) [8]
	tg.food(6, 0)       // food at (6,0)
	tg.eat()            // (5,0) -> (6,0)::(3,5) [9]
	tg.food(7, 0)       // food at (7,0)
	tg.eat()            // (6,0) -> (7,0)::(3,5) [10]
	tg.food(8, 0)       // food at (8,0)
	tg.eat()            // (7,0) -> (8,0)::(3,5) [11]
	tg.food(9, 9)       // food at (9,0)
	tg.turn(game.Up)    // turn up
	tg.move(2)          // (8,0) -> (8,2)::(3,3)
	tg.turn(game.Left)  // turn left
	tg.move(2)          // (8,2) -> (6,2)::(3,1)
	tg.turn(game.Down)  // turn down
	tg.move(1)          // (6,2) -> (6,1)::(3,0)

	// The next step should collide with itself on the next step (6,0)
	tg.collide()
}

func Test_GamePlayBoundaryCollision(t *testing.T) {
	tg := testingGame(t)

	tg.move(4)          // (5,5) -> (5,9)
	tg.turn(game.Right) // turn right
	tg.move(4)          // (5,9) -> (9,9)
	tg.turn(game.Down)  // turn down
	tg.move(9)          // (9,9) -> (9,0)

	// Expect a boundary collision on the next step
	tg.boundary()
}

/**
 * Some tools used to make gmae path testing more straightforward
 */

// TestingGame a verbose game with testing elements mixed in, and methods for executing
// actions with logged results
type TestingGame struct {
	game *game.Game
	t    *testing.T
}

type tickAction struct {
	Result   game.TickResult
	HasError bool
}
type foodAction struct {
	Food game.Point
}
type turnAction struct {
	Dir game.Vector
}

// Move n spaces
func (tg TestingGame) move(n int) {
	a := tickAction{Result: game.TickResult{Moved: true}}
	for i := 0; i < n; i++ {
		processTick(a, tg.game, tg.t)
	}
}

// Move and expect an eat
func (tg TestingGame) eat() {
	processTick(tickAction{Result: game.TickResult{AteFood: true, Grew: true}}, tg.game, tg.t)
}

// Perform a turn
func (tg TestingGame) turn(d game.Vector) {
	processTurn(turnAction{Dir: d}, tg.game, tg.t)
}

// Set a new food point
func (tg TestingGame) food(x, y int) {
	f := game.Point{X: x, Y: y}
	processFood(foodAction{Food: f}, tg.game, tg.t)
}

// Move and expect a collision
func (tg TestingGame) collide() {
	processTick(tickAction{Result: game.TickResult{SnakeCollision: true}, HasError: true}, tg.game, tg.t)
}

// Move and expect a collision
func (tg TestingGame) boundary() {
	processTick(tickAction{Result: game.TickResult{BoundaryCollision: true}, HasError: true}, tg.game, tg.t)
}

// Process test actions as individual logged actions

// Process a slice of actions
func processActions(g *game.Game, t *testing.T, as []interface{}) {
	for _, a := range as {
		switch a.(type) {
		case turnAction:
			processTurn(a.(turnAction), g, t)
		case tickAction:
			processTick(a.(tickAction), g, t)
		case foodAction:
			processFood(a.(foodAction), g, t)
		default:
			t.Errorf("Could not process action of unknown type")
		}
	}
}

// use a standard function for turning
func processTurn(a turnAction, g *game.Game, t *testing.T) error {
	t.Logf("TURN: [Direction: %s -> %s]", g.Facing(), a.Dir)
	g.Turn(a.Dir)
	return nil
}

// use a standard function for setting new food point
func processFood(a foodAction, g *game.Game, t *testing.T) error {
	t.Logf("FOOD: [Position: %s]", a.Food)
	g.SetFood(a.Food)
	return nil
}

// Use a standard function for ticking the game, and logging the game status
func processTick(a tickAction, g *game.Game, t *testing.T) (game.TickResult, error) {
	res, err := g.Tick()
	aRes := a.Result

	resString := "unknown"
	if a.HasError && aRes.BoundaryCollision {
		resString = "EXPECTED BOUNDARY COLLISION"
	} else if a.HasError && aRes.SnakeCollision {
		resString = "EXPECTED SNAKE COLLISION"
	} else if a.HasError {
		resString = "EXPECTED ERROR"
	} else if err != nil {
		resString = "ERROR"
	} else if res.Grew {
		resString = "ATE/GREW"
	} else if res.Moved {
		resString = "MOVED"
	}

	if a.HasError && err == nil {
		t.Errorf("Tick action expected an error, but did not get one")
	} else if err != nil && !a.HasError {
		t.Errorf("Tick action produced an unexpected error : %s", err.Error())
	}

	if aRes.BoundaryCollision != res.BoundaryCollision {
		if aRes.BoundaryCollision {
			t.Errorf("Expected boundary collision not detected")
		} else {
			t.Errorf("Unexpected boundary collision detected")
		}
	}
	if aRes.SnakeCollision != res.SnakeCollision {
		if aRes.SnakeCollision {
			t.Errorf("Expected snake collision not detected")
		} else {
			t.Errorf("Unexpected snake collision detected")
		}
	}
	if aRes.AteFood != res.AteFood {
		if aRes.AteFood {
			t.Errorf("Expected eat event not detected")
		} else {
			t.Errorf("Unexpected eat event detected")
		}
	}
	if aRes.Grew != res.Grew {
		if aRes.Grew {
			t.Errorf("Expected grow not detected")
		} else {
			t.Errorf("Unexpected grow detected")
		}
	}
	if aRes.Moved != res.Moved {
		if aRes.Moved {
			t.Errorf("Expected move not detected")
		} else {
			t.Errorf("Unexpected move collision detected")
		}
	}

	t.Logf("TICK: %s [Direction: %s][Snake: [%v] %s]", resString, g.Facing(), g.Length(), g.Head())
	return res, err
}
