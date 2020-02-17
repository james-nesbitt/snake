# Game

The data structures and business logic for a game of snake

## Space

Spacial elements used as the foundation for the interactions of the game.
These are specifically Cartesian elements, principly in 2D.  We try to keep all
dimension elements here so that the space can be considered an abstract.
Concepts like collision, orthogonality and equivalence are kept here.

Point/Vector : single set of Cartesian coordinate.  A point in space or a path
  from origin to the point as a vector.

  There are a set of Unit Vectors that can be used to indicate turns for snake
  vectoring

Grid : a vector which indicates the dimensions of the game space.  It can also
  determine if a given point is inside of outside of the grid.

## Snake

Game data structures used to represent elements of the game.

Food : a point where the food is

Snake:
  Element : a linked list node for a piece of the snake
  Snake : The head of the snake and some functionality for moving/growing and
    detecting if a point is in the snake.

    There isn't any full cycle detection in the snake, as the snake is really a
    queue where only the new head needs to be cycle tested.

## Game

A game combines spacial elements and Snake elements into a single set, which can
be used to manage a full game.

The game is also where snake collisions (self and boundary) are checked.  A snake
collision could be detected in the snake, but a snake-boundary collision cannot
and a snake-food collision (EAT) cannot, so we might as well do all collision
event detection here.
