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
  * $A-A,B-\emptyset : \varepsilon$
  * <img src="http://latex.codecogs.com/gif.latex?\frac{\partial J}{\partial \theta_k^{(j)}}=\sum_{i:r(i,j)=1}{\big((\theta^{(j)})^Tx^{(i)}-y^{(i,j)}\big)x_k^{(i)}}+\lambda \theta_k^{(j)}" />
