```
.     .       .  .   . .   .   . .    +  .
  .     .  :     .    .. :. .___---------___.
       .  .   .    .  :.:. _".^ .^ ^.  '.. :"-_. .
    .  :       .  .  .:../:            . .^  :.:\.
        .   . :: +. :.:/: .   .    .        . . .:\
 .  :    .     . _ :::/:               .  ^ .  . .:\
  .. . .   . - : :.:./.                        .  .:\
  .      .     . :..|:                    .  .  ^. .:|
    .       . : : ..||        .                . . !:|
  .     . . . ::. ::\(                           . :)/
 .   .     : . : .:.|. ######              .#######::|
  :.. .  :-  : .:  ::|.#######           ..########:|
 .  .  .  ..  .  .. :\ ########          :######## :/
  .        .+ :: : -.:\ ########       . ########.:/
    .  .+   . . . . :.:\. #######       #######..:/
      :: . . . . ::.:..:.\           .   .   ..:/
   .   .   .  .. :  -::::.\.       | |     . .:/
      .  :  .  .  .-:.":.::.\             ..:/
 .      -.   . . . .: .:::.:.\.           .:/
.   .   .  :      : ....::_:..:\   ___.  :/
   .   .  .   .:. .. .  .: :.:.:\       :/
     +   .   .   : . ::. :.:. .:.|\  .:/|
     .         +   .  .  ...:: ..|  --.:|
.      . . .   .  .  . ... :..:.."(  ..)"
 .   .       .      :  .   .: ::/  .  .::\
```

#Alien Invation

---

### Assumptions
- Input file is valid. If `Foo north=Bar` is true then `Bar south=Foo` must hold true.
- City names do not contain special characters.
- Directions are `north`, `south`, `east` and `west`. 
- Roads are bidirectional.
- Number of aliens can be more than the number of cities and at most 10000.
- Dataset is small. All the cities and the roads can be stored in memory of a single machine.
- If an alien is trapped and can not make a move, we consider it made a move after one epoch has passed. Otherwise, we can end up in a livelock.
- More than one alien can converge to a city. They destroy each other and the city.

### Algorithmic Assumptions
- __For running the simulation__, it only matters if two city are connected. The orientation (north, south...) is not important.
- The program will run on a machine where word size is at least 8 bits. Which is true for almost all modern machines. 

### How to Run
from the root of the project (the folder that contains go.mod file) `go mod tidy && go run main.go 3`
where 3 is the number of aliens. Number of aliens can be in range [0, 10000]

This will read from a file named `data.txt` located at the root of the project.

**For test:** From the root of the project `go test -v ./...`

### Philosophy
- Keep things simple. Explicitly! I used `int` in place of `int32` because what are the ods that we run this application on a machine where word size is not 8? And does that affect the correctness? I do understand that for cross platforms it's best to fix the size.
- I've decided to treat the termination condition `iteration number reached limit[10000]` as an error. Because for a lot of case the simulation can go on for an arbitrary long time. And we can be __terminating__ the simulation __prematurely__. It can happen that on the 10001 th tern, two aliens collide. So, when the simulation reached the 10000th step we terminate prematurely with and error.
- I've used one external library. It could be avoided. But by the time I realised I am using a 3rd party library which could be avoided, I've already written some tests with it. So I just went with the flow. No philosophy here :)  
