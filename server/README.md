# Server

A snake Server is a chan based struct which can be used to interact with a snake
game.  It is meant to be a single point of interaction for a game UI to use for
running a game.

The server provides outgoing chans for event response:
1. New food is needed (the snake ate the food)
2. The snake hit the boundary - game end
3. The snake hit itself - game end

The server provides incoming chans for game progress and snake control:
1. a game clock tick
2. a snake turn chan for controlling the snake .

To get UI info, the Server Game can be directly read from to get the grid
dimensions, the snake points/facing-direction and the food position.

The packages server constructor expects you to create your game instance outside
of the Server. This means that you will want to set your grid size, you initial
snake position, and the first food position before creating the server.

## Needs Food

The signal for needing new food is a chan of a chan.

```
  NeedsFood chan chan game.Point
```

The idea is that a signal will go out that food is needed.  The signal is a chan
of a game point, which should be used to send the new food position.

The Server will block Tick and Turn events until a new food position is placed.

This was not a technical requirement, but was done to simplify the game loop, so
that we did not have to detect ticks where we need to wait for food, and so that
we did not have to correlate new food position events with requests on a single
chan.

I think this makes that game loop more stable, and interactions more clear, at
the cost of a responsibility of providing a new food position before you can
tick.
