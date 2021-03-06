# Spaloosh!
[Play now on itch.io](https://atolverderben.itch.io/spaloosh)
![Spaloosh Title Screen](screenshots/title.png "Spaloosh! Title Screen")


## About
Spaloosh is a take on the traditional battleship game using a monster theme instead of battleships. Originally created for the [Maximum Monster Month](https://itch.io/jam/maximum-monster-month) game jam on [itch.io](https://itch.io). It is written in Go using [Ebiten](https://github.com/hajimehoshi/ebiten) for a 2D framework. All "artwork" was drawn by myself. The idea was to challenge myself to make a "full" game from scratch all by myself, that included using no established game engines.

## Features
* Two different game modes
  * Race against the clock to sink 3 monsters with limited bombs. Includes 2 difficulty levels
  * Battle against the AI in traditional battleship style matches.
* Online Multiplayer Matches
  * Direct IP connection
  * Server Rooms available (players may host their own servers with the provided code)
* Multiple Characters to choose from (all need a backstory though that is in progress)
* Cross Platform
  * Since the entire game is written using the Go language it can easily be cross compiled to all major operating systems

## TODO
* Code Cleanup
  * Game Jams lead to lots of messy code, such as copy paste ugliness
* Improved AI
  * Due to limitations in time for the jam the AI isn't the most advanced
  * I had started work on more advanced features I would like to continue
* Better Online Support
  * Currently the online is mostly centered on ad hoc play with direct IP connections
  * Would like to add better server implementations where the server is the authority
* Possible character backgrounds of why all these random monsters are playing this game!
* Potentially add more board games to the mix (this might be a new project though)


<br/>
<br/>

  ![Gameplay start example](screenshots/startgame.gif "Game Start Example")
