package server_test

import (
	"context"
	"github.com/james-nesbitt/snake/game"
	"github.com/james-nesbitt/snake/server"
	"testing"
	"time"
)

// Test the NeedsFood handler
func Test_NeedsFoodBase(t *testing.T) {
	nfc := make(chan chan game.Point)
	nf := make(chan game.Point)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ef := game.Point{X: 3, Y: 4}
	mf := &NeedsFood_Mock{Food: ef}

	go server.NeedFoodHandler(mf, nfc, ctx)

	nfc <- nf
	gf := <-nf
	if !gf.Equals(ef) {
		t.Errorf("Needs Food handling function did not return expected food")
	}
}

// Test the moving NeedsFood handler
func Test_NeedsFoodMove(t *testing.T) {
	nfc := make(chan chan game.Point)
	nf := make(chan game.Point)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mv := game.Vector{X: 1, Y: 2}
	f := game.Point{X: 3, Y: 2}
	g, err := game.AutoGame(game.Vector{X: 10, Y: 10}, f)
	if err != nil {
		t.Errorf("Game construction error: %s", err)
	}

	mf := server.NewMakeFood_Move(&g, mv)
	go server.NeedFoodHandler(mf, nfc, ctx)

	if gf, err := g.Food(); err != nil {
		t.Errorf("Could not get food position: %s", err)
	} else if !f.Equals(gf) {
		t.Errorf("Food is in a weird place ")
	}

	for i := 0; i < 5; i++ {
		nfc <- nf
		gf := <-nf

		g.SetFood(gf)
		t.Logf("Food moved to %s", gf)
		f = f.Move(mv)
		if !gf.Equals(f) {
			t.Errorf("Move Food handler did not move the food as expected: %s / %s", gf, f)
		} else {
			t.Log("Good move")
		}
	}
}

// Test slice based NeedsFood handler
func Test_NeedsFoodSlice(t *testing.T) {
	nfc := make(chan chan game.Point)
	nf := make(chan game.Point)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	f := game.Point{X: 3, Y: 2}
	g, err := game.AutoGame(game.Vector{X: 10, Y: 10}, f)
	if err != nil {
		t.Errorf("Game construction error: %s", err)
	}

	fs := []game.Point{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 0, Y: 8},
		{X: 7, Y: 7},
		{X: 3, Y: 5},
		{X: 3, Y: 7},
	}
	mf := server.NewMakeFood_Slice(fs)
	go server.NeedFoodHandler(mf, nfc, ctx)

	if gf, err := g.Food(); err != nil {
		t.Errorf("Could not get food position: %s", err)
	} else if !f.Equals(gf) {
		t.Errorf("Food is in a weird place ")
	}

	for _, ef := range fs {
		nfc <- nf
		gf := <-nf

		g.SetFood(gf)
		t.Logf("Food moved to %s", gf)
		if !gf.Equals(ef) {
			t.Errorf("Move Food handler did not move the food as expected: %s / %s", gf, ef)
		} else {
			t.Log("Good move")
		}
	}
}

// Movk for NeedsFood - it just returns whatever point if was filled with
type NeedsFood_Mock struct {
	Food game.Point
}

func (nfm NeedsFood_Mock) NextFood() game.Point {
	return nfm.Food
}
