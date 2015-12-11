# reo-learn

This is a project using active learning technique to extract a reo connector model.

All the codes are written under Google Go.

# project structure

* lib
  * reo - this library includes the definition of Ports, basic Reo Channels, etc.
  * sul - 
  * learn

* example
  * fifo
  * alternator

## several questions

* reo's semantics when a write operation from components is blocked?

## TODO

* tree optimization in MQuery
* debug of StepExec
* solve the errors in table
  * ![](http://latex.codecogs.com/gif.latex?A-A,B-\\emptyset:\\varepsilon)
  * A,B-A,B-A | C:B, C:B,  C:A,  C:A,  ϵ
    * bug in merger channel leads to this bug

## logs

- **Dec 11 2015** problems on Merger are solved
