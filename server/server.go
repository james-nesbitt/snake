package server

import (
	"context"
	"github.com/james-nesbitt/snake/game"
	"log"
)

func NewServer(g *game.Game) Server {
	tk := make(chan int)
	tn := make(chan game.Vector)
	nf := make(chan chan game.Point)
	bc := make(chan error)
	sc := make(chan error)

	return Server{Game: g, Tick: tk, Turn: tn, NeedsFood: nf, BoundaryCollision: bc, SnakeCollision: sc}
}

/**
 * A server object which can be used to interact with a game as a backend
 */
type Server struct {
	Game *game.Game
	// Incoming instructions
	Tick chan int         // Game tick (step)
	Turn chan game.Vector // Snake turn
	// Outgoing info
	NeedsFood chan chan game.Point // Food was eaten (new food needed on the passed chan)

	// Outgoing errors
	BoundaryCollision chan error
	SnakeCollision    chan error
}

func (s *Server) Start(ctx context.Context) {
	log.Printf("START SNAKE SERVER")

	/**
	 * Main event loop
	 *
	 * Here we wait for the outer signals that run the game.
	 *
	 * The loops is signal (chan) based, reacting only to signals. The construction
	 * context can be used to externally shut down the event loop.
	 *
	 * INCOMING SIGNALS
	 * 1. TICK -> a game clock tick.  Done as a signal so that a  Game consumer can
	 *            regulate the game
	 * 2. TURN -> a snake turn even. This can happen at any time
	 * 3. FOOD -> new food placement.  This is EXPECTED to occur only after the
	 *            outgoing NEEDFOOD signal is sent, but really it could happen
	 *  @TODO  perhaps the NEEDFOOD chan needs to be reworked to be isolated.
	 *
	 * OUTGOING SIGNALS
	 *
	 * 4. NEEDFOOD <- New food is needed.  An external algorithm should be applied
	 *           which decides where to put food.
	 *
	 * ERROR SIGNALS (OUTGOING)
	 *
	 * 5. BOUNDARY -> The snake has collided with the grid boundary
	 * 6. SNAKECOLLISION -> The snake has collided with itself
	 *
	 */

	for {
		select {
		case <-ctx.Done():
			log.Printf("STOP REQUESTED")
			s.stop()
			return
		case _ = <-s.Tick:
			res, err := s.Game.Tick()

			if err != nil {
				log.Printf("TICK: ERROR [Snake: %s]", s.Game.Head())
				if res.BoundaryCollision {
					s.BoundaryCollision <- err
				}
				if res.SnakeCollision {
					s.SnakeCollision <- err
				}

				s.stop()

				return
			} else if res.Grew {
				log.Printf("TICK: GREW [Dir: %s][Snake: %s]", s.Game.Facing(), s.Game.Head())
			} else if res.Moved {
				if f, err := s.Game.Food(); err != nil {
					log.Printf("TICK: MOVED [Dir: %s][Food: %s][Snake: %s]", s.Game.Facing(), "NONE", s.Game.Head())
				} else {
					log.Printf("TICK: MOVED [Dir: %s][Food: %s][Snake: %s]", s.Game.Facing(), f, s.Game.Head())
				}
			}

			if res.AteFood {
				// originally we played with separation of the NeedsFood and Food chans
				// but it required validation on the tick level and caused an issue with
				// closed channels if making food happens after closing the outer context

				log.Printf("ATE: Asking for new food point")
				fc := make(chan game.Point) // New food chan, to receive a new food point
				s.NeedsFood <- fc           // send out a signal that we need new food
				food := <-fc                // receive new food position from the sent chan
				close(fc)                   // prevent subsequent sends to the chan
				s.Game.SetFood(food)        // place the food
				log.Printf("FOOD: New food created at %s", food)
			}

		case dir := <-s.Turn:
			log.Printf("TURNED: %s -> %s ", s.Game.Facing(), dir)
			s.Game.Turn(dir)
		}
	}
}

// Stop the Server
func (s *Server) stop() {
	close(s.Tick)
	close(s.Turn)
	close(s.NeedsFood)
	close(s.BoundaryCollision)
	close(s.SnakeCollision)
	log.Printf("STOPPED SNAKE SERVER")
}
