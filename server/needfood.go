package server

import (
	"context"
	"github.com/james-nesbitt/snake/game"
	"log"
	"math/rand"
)

/**
 * Some helper functionality for Responding to the NeedsFood signal
 */

// Something that can MakeFood points
type MakeFood interface {
	NextFood() game.Point
}

// A handler function that can be put in charge of making Food, when asked for
// @USAGE use this as a subroutine for responding to a NeedsFood chan
func NeedFoodHandler(mf MakeFood, nf chan chan game.Point, ctx context.Context) {
	log.Printf("Starting to listen for NeedFood events")
	for {
		select {
		case fc := <-nf:
			log.Printf("Received request for new Food location")
			fc <- mf.NextFood()
			log.Printf("Sent new food location")

			// @NOTE originally we closed the channel, but most implementations will
			//       reuse that passed food chan, and we should let them close it.
		case <-ctx.Done():
			return
		}
	}
}

func NewMakeFood_Move(g *game.Game, m game.Vector) MakeFood {
	return &MakeFood_Move{g: g, m: m}
}

// Move the food in a Vector every time
type MakeFood_Move struct {
	g *game.Game
	m game.Vector
}

func (mf *MakeFood_Move) NextFood() game.Point {
	hp, _ := mf.g.Food()
	np := hp.Move(mf.m)
	log.Printf("Moved food from %s to %s", hp, np)
	return np
}

func NewMakeFood_Random(g *game.Game) MakeFood {
	return &MakeFood_Random{g: g}
}

// Return the next point from a source slice
type MakeFood_Random struct {
	g *game.Game
}

func (mf *MakeFood_Random) NextFood() game.Point {
	var f game.Point // Use this to hold food points while we try to find one
	sz := mf.g.Size()
	for {
		f := game.Point{X: rand.Intn(sz.X), Y: rand.Intn(sz.Y)}
		if !mf.g.Head().FindPoint(f) {
			break
		}
	}
	return f
}

func NewMakeFood_Slice(ps []game.Point) MakeFood {
	return &MakeFood_Slice{Points: ps}
}

// Return the next point from a source slice
type MakeFood_Slice struct {
	Points []game.Point
}

func (mf *MakeFood_Slice) NextFood() game.Point {
	var f game.Point
	f, mf.Points = mf.Points[0], mf.Points[1:]
	return f
}
