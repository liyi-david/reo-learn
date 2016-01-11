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
  * 2-buffer (mainly used to test the counter-example processing algorithm)

## TODO
- abstract
- full paper
- there's problem in 2-buffer. in the second round, obstables cannot be displayed properly *partly solved on my laptop*
- bugs found in new example

## logs

- **Dec 11 2015** problems on Merger are solved
- **Dec 17 2015** rewrite the expandLp(). with the new expand function, the performance has been improved greatly.
- **Dec 22 2015** *Equivalence Query* online now. Also the counter-example analyzing function has passed the test
- **Dec 25 2015** timer channels and a corresponding example has been online
- **Jan 08 2016** redundance added in sequence simulation
- **Jan 09 2016** time-comsumption analysis is added in both sul and learn modules
- **Jan 09 2016** new tree-optimization tactic
- **Jan 11 2016** new example showing how to use reo package for concurrent programming
