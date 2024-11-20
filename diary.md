## Development Diary

- 18/11/2024 Day 1
  Started collecting tools and ideas. Discovered [a project](https://github.com/eliukblau/pixterm) that takes a picture and displays it on terminal, can even emit Go code; cleaned up the HW logo and made some experiments on image processing

- 19/11/2024 Day 2
  Study [BubbleTea Framework](https://github.com/charmbracelet/bubbletea) and the fundamentals of the [ELM Architecture](https://guide.elm-lang.org/architecture/); implemented a simple animation system with constant-frame update, with a tick timer that send a custom "Frame" message every Nth of a second

- 20/11/2024 Day 3
  Picked [maze for programmers](http://www.mazesforprogrammers.com/) from my reading list and playing with an implementation that allows double vertical resolution by using half-block unicode characters; also fiddling around with [lipgloss](https://github.com/charmbracelet/lipgloss) styling

- 21/11/2024 Day 4
  Re-used part of the maze code to implement a simple game. I tried to separate code in two different packages and now it runs without any timer, changing state (and using system resources) only on user input events. Now each maze cell is two characters wide, but in the model is still represented as one byte.

  

