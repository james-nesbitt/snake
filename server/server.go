package server

/**
 * A snake server. The component which can be used to handle snake mechanics
 * for a front end, using channels for interaction.
 *
 * The Server is a struct which provides incoming and outgoing channels:
 *   tick : a clock tick in the snake game (incoming)
 *   turn : a snake direction turn event (incoming)
 *   needs-food : new food placement is needed (food was eaten)
 *   collision-boundary : the snake ran into the grid boundary (outgoing)
 *   collision-snake : the snale ran into itself (outgoing)
 *
 * The server must be "Start"ed before interacting with the channels, which
 * needs a context that can be used to kill the Server game.
 *
 * @TODO the Start method has a context which can be used to stop it, so the
 *    Stop method should probably be removed
 */

import (
	"context"
	"github.com/james-nesbitt/snake/game"
	"log"
)

// NewServer server constructor.  Don't forget to Start before using it
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
 *
 * @TODO chans should be strict about incoming and outgoing.
 */
type Server struct {
	Game *game.Game

	// Incoming instructions
	Tick chan int         // Game tick (step) trigger
	Turn chan game.Vector // Snake turn trigger

	// Outgoing info
	NeedsFood chan chan game.Point // Food was eaten (new food needed on the passed chan)

	// Outgoing errors
	BoundaryCollision chan error
	SnakeCollision    chan error
}

// Start the server running a game by open all channels and listening on then in
// a game loop
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
