package server_test

import (
	"context"
	"github.com/james-nesbitt/snake/game"
	"github.com/james-nesbitt/snake/server"
	"testing"
	"time"
)

var (
	testTick = time.Millisecond * 5 // a tick that we can use for testing, so
	                                // that we can turn and place food as well
)

// Use a clock ticker and do some simple snake game, simple checking
func Test_BasicTimedServer(t *testing.T) {
	ticker := time.NewTicker(testTick)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	g, err := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 5, Y: 7})
	if err != nil {
		t.Errorf("Game construction error: %s", err)
	}
	s := server.NewServer(&g)

	go logErrorChan(s.BoundaryCollision, t)
	go logErrorChan(s.SnakeCollision, t)

	// place new food two points ahead of the snake
	go server.NeedFoodHandler(server.NewMakeFood_Move(&g, game.Vector{Y: 2}), s.NeedsFood, ctx)

	go s.Start(ctx)

	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			s.Tick <- i
		case <-ctx.Done():
			t.Errorf("Safe timeout expired")
			return
		}
	}

	cancel()
}

// Test some game snake wandering.  Move and Grow
func Test_SnakeWanderServer(t *testing.T) {
	ticker := time.NewTicker(testTick)
	defer ticker.Stop()
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	g, _ := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 5, Y: 7})
	s := server.NewServer(&g)

	go logErrorChan(s.BoundaryCollision, t)
	go logErrorChan(s.SnakeCollision, t)

	// Move the food on a planned set of points
	m := []game.Point{
		{X: 1, Y: 1},
		{X: 7, Y: 8},
		{X: 2, Y: 9},
		{X: 5, Y: 0},
	}
	go server.NeedFoodHandler(server.NewMakeFood_Slice(m), s.NeedsFood, ctx)

	go s.Start(ctx)

	for i := 0; i < 4; i++ {
		select {
		case <-ticker.C:
			s.Tick <- i
		case <-ctx.Done():
			t.Errorf("Safe timeout expired")
			return
		}
	}
}

// Test some game snake wandering.  Move and Grow along a significant path
func Test_SnakePlayAGame(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	g, _ := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 5, Y: 7})
	s := server.NewServer(&g)

	go logErrorChan(s.BoundaryCollision, t)
	go logErrorChan(s.SnakeCollision, t)

	// Move the food on a planned set of points
	m := []game.Point{
		{X: 7, Y: 8},
		{X: 2, Y: 9},
		{X: 5, Y: 0},
		{X: 3, Y: 2},
		{X: 1, Y: 1},
		{X: 4, Y: 8},
		{X: 8, Y: 7},
	}
	go server.NeedFoodHandler(server.NewMakeFood_Slice(m), s.NeedsFood, ctx)

	go s.Start(ctx)

	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Turn <- game.Right
	s.Tick <- 2
	s.Tick <- 2
	s.Turn <- game.Up
	s.Tick <- 3
	s.Turn <- game.Left
	s.Tick <- 4
	s.Tick <- 4
	s.Tick <- 4
	s.Tick <- 4
	s.Tick <- 4
	s.Turn <- game.Down
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Tick <- 5
	s.Turn <- game.Right
	s.Tick <- 6
	s.Tick <- 6
	s.Tick <- 6
	s.Tick <- 6
	s.Turn <- game.Up
	s.Tick <- 7
	s.Tick <- 7
	s.Turn <- game.Left
	s.Tick <- 8
	s.Tick <- 8
	s.Tick <- 8
	s.Tick <- 8
	s.Tick <- 8
	s.Turn <- game.Down
	s.Tick <- 9
	s.Tick <- 9
}

// Test Move and Grow and end up in a self-collision
func Test_SnakeSelfCollide(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	g, _ := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 5, Y: 6})
	s := server.NewServer(&g)

	go logErrorChan(s.BoundaryCollision, t)

	// Listen for unexpected errors until we say stop
	stop := make(chan bool)
	go func(err chan error, stop chan bool, t *testing.T) {
		select {
		case rec := <-err:
			t.Error("Error received: ", rec)
		case <-stop:
			t.Log("Disabling Snake collision catch")
			return
		}
	}(s.SnakeCollision, stop, t)

	// Move the food on a planned set of points which gives us a fast chance to grow
	m := []game.Point{
		{X: 5, Y: 7},
		{X: 5, Y: 8},
		{X: 5, Y: 9},
		{X: 5, Y: 10},
		{X: 6, Y: 10},
		{X: 7, Y: 10},
		{X: 8, Y: 10},
		{X: 1, Y: 11}, // this one we never eat but we need it to prevent index issues
	}
	go server.NeedFoodHandler(server.NewMakeFood_Slice(m), s.NeedsFood, ctx)

	go s.Start(ctx)

	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Turn <- game.Right
	s.Tick <- 2
	s.Tick <- 2
	s.Turn <- game.Down
	s.Tick <- 3
	s.Turn <- game.Left
	s.Tick <- 4

	stop <- true // the SnakeBoundary Chan now has no listener

	s.Tick <- 4 // Should cause a snake collision

	giveup, _ := context.WithTimeout(context.Background(), 3*time.Second)
	select {
	case <-giveup.Done():
		t.Errorf("Failed to receive expected snake collision error on chan")
	case err := <-s.SnakeCollision:
		t.Logf("Received expected error on snake collision error chan : %s", err)
	}
}

// Test some game snake wandering.  Move and Grow
func Test_SnakeBoundaryCollide(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	g, _ := game.AutoGame(game.Vector{X: 10, Y: 10}, game.Point{X: 5, Y: 6})
	s := server.NewServer(&g)

	// Listen for unexpected errors until we say stop
	stop := make(chan bool)
	go func(err chan error, stop chan bool, t *testing.T) {
		select {
		case rec := <-err:
			t.Error("Error received: ", rec)
		case <-stop:
			t.Log("Disabling Boundary collision catch")
			return
		}
	}(s.BoundaryCollision, stop, t)

	go logErrorChan(s.SnakeCollision, t)

	// We won't eat food here, so we can ignore food generation
	go server.NeedFoodHandler(NeedsFood_Mock{Food: game.Point{X: 1, Y: 2}}, s.NeedsFood, ctx)

	go s.Start(ctx)

	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Tick <- 1
	s.Turn <- game.Right
	s.Tick <- 2
	s.Tick <- 2
	s.Tick <- 2
	s.Tick <- 2
	s.Tick <- 2

	stop <- true // the SnakeBoundary Chan now has no listener

	s.Tick <- 4 // Should cause a snake collision

	giveup, _ := context.WithTimeout(context.Background(), 3*time.Second)
	select {
	case <-giveup.Done():
		t.Errorf("Failed to receive expected boundary collision error on chan")
	case err := <-s.BoundaryCollision:
		t.Logf("Received expected error on boundary collision error chan : %s", err)
	}
}

// Just log errors if they come in - these should be unexpected errors that you
// don't want to catch yourself
func logErrorChan(err chan error, t *testing.T) {
	if rec := <-err; rec != nil {
		t.Error("Error received: ", rec)
	}
}
