### Assumptions
- Number of aliens can be at most 10000.
- Dataset is small. All the cities and the roads can be stored in memory of a single machine.
- Roads are bidirectional.
- If an alien is trapped and can not make a move, we consider it made a move after one epoch has passed. Otherwise, we can end up in a livelock.
- More than one alien can converge to a city. They destroy each other and the city.

### Algorithmic Assumptions
- It only matters if two city are connected. The connection orientation is not important.