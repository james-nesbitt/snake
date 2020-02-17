package main

import (
	"context"
  "fmt"
	"github.com/james-nesbitt/snake/game"
	"github.com/james-nesbitt/snake/server"
	"github.com/jroimartin/gocui"
	"log"
	"time"
)

var (
	grid       = game.Vector{X: 50, Y: 50}
	foodStart  = game.Point{X: 0, Y: 0}
	tickPeriod = time.Second

	i int // game ticks
  t *time.Ticker
  ctx context.Context
  cancel func()
  g game.Game
  err error
  s server.Server
  mf server.MakeFood
	h []string

  gu *gocui.Gui
	hv *gocui.View  // History view
	gv *gocui.View // Game/Grid view
	lv *gocui.View // Log view
)

func init() {
	i = 0
	t = time.NewTicker(tickPeriod)
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

	g, err = game.AutoGame(grid, foodStart)
	if err != nil {
		log.Printf("Game construction error: %s", err)
	}

	s = server.NewServer(&g)
  mf = server.NewMakeFood_Random(&g)

	h = []string{}
}


func main() {
	var err error
	gu, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer gu.Close()

	gu.SetManagerFunc(draw)

	if err := gu.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gu.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	// select {
	// case <- t.C:
	// 	s.Tick <- i
	// 	i += 1
	// 	gu.Update(func(gu *gocui.Gui) error {
	// 		updateGrid()
	// 		return nil
	// 	})
	// }
}

func draw(g *gocui.Gui) error {
	var err error
	if hv, err = g.SetView("history", 0, 0, 29, grid.Y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if gv, err = g.SetView("grid", 30, 0, 30+grid.X, grid.Y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if lv, err = g.SetView("log", 31+grid.X, 0, 81+grid.X, grid.Y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	updateHistory()
	updateGrid()

	return nil
}

func updateHistory() {
	hv.Clear()
	for i, s := range h {
		fmt.Fprintln(hv, "%3d) %s", i, s)
	}
}
func updateGrid() {
	gv.Clear()

	if f, err := s.Game.Food(); err != nil {
		log.Printf("%s", err)
	} else {
		gv.SetCursor(f.X, f.Y)
		gv.Write([]byte("F"))
	}

	for ps := s.Game.Head(); ps != nil; ps = ps.Next() {
		p := ps.Point()
		gv.SetCursor(p.X, p.Y)
		gv.Write([]byte("S"))
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// func main() {
// 	defer t.Stop()
//
// 	// place new food in random places
// 	go server.NeedFoodHandler(mf, s.NeedsFood, ctx)
//
// 	// Get input from stdin, for turns and stops
// 	in := keypress()
//
// 	// start the game server
// 	// s.Start(ctx)
//
//   for i:= 0; i>=0; i++{
//   	select {
//   	case r := <-in:
//   		switch r {
//   		// case 'w':
//   		// 	s.Turn <- game.Up
//   		// case 's':
//   		// 	s.Turn <- game.Down
//   		// case 'a':
//   		// 	s.Turn <- game.Left
//   		// case 'd':
//   		// 	s.Turn <- game.Right
//   		case 'q':
//   			return
//   			cancel()
// 				t.Stop()
//   			log.Print("Ending game")
// 			default:
// 				log.Printf("Keypress %s", r)
//   		}
//   	// case <-t.C:
//   	// 	s.Tick <- i
//   	// case err := <-s.SnakeCollision:
//   	// 	log.Print("Snake collided with itself: %s", err)
//   	// 	cancel()
//   	// 	return
//   	// case err := <-s.BoundaryCollision:
//   	// 	log.Print("Snake collided with the boundary: %s", err)
//   	// 	cancel()
//   	// 	return
//   	}
//   }
// }
