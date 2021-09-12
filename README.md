### Assumptions
- Input file is valid. If `Foo north=Bar` is true then `Bar south=Foo` must hold true.
- Number of aliens can be at most 10000.
- Dataset is small. All the cities and the roads can be stored in memory of a single machine.
- Roads are bidirectional.
- If an alien is trapped and can not make a move, we consider it made a move after one epoch has passed. Otherwise, we can end up in a livelock.
- More than one alien can converge to a city. They destroy each other and the city.

### Algorithmic Assumptions
- It only matters if two city are connected. The connection orientation is not important.


### Phylosophy
- Keep things simple. Explicitly! I used `int` in place of `int32` because what are the ods that we run this application on a machine where int word size is not 8? And does that affect the correctness? I do understand the machine architecture and know when to fix the byte size for int.
- 