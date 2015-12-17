# reo-learn

This is a project using active learning technique to extract a reo connector model.

All the codes are written under Google Go.

# project structure

* lib
  * reo - Ports, basic Reo Channels, etc.
  * sul - SUL(System Under Learn), Oracles, Membership Query, and *Equivalence Query*
  * learn

* example
  * fifo
  * alternator

## several questions

* reo's semantics when a write operation from components is blocked?

## TODO
- write the paper ...
- realize the EQuery function

## logs

- **Dec 11 2015** problems on Merger are solved
- **Dec 17 2015** rewrite the expandLp(). with the new expand function, the performance has been improved greatly.
