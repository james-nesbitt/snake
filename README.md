# Snake

The snake game, written in go.

I wrote this as programming challenge as I was getting back into go development
and I wanted to make sure that I hadn't lost my touch.

The code is separated into the following parts (go modules):

1. Game : native elements of the snake game:
   a. spacial components such as grids, points, vectors (snake direction)
   b. snake elements such as snake pieces, the snake and food
2. Server : an interactive server object which uses channels for interaction.
3. UIs :
   a. a screen ui
   b. (COMING SOON) an html ui

## Building

Currently this is just demo and doesn't have building instructions.  When the UI
is developed, then build instructions will center around that, so the build
instructions may be found in the particular UI approach that you are looking for
