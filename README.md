### PISMO
Implementation of PhyloSor similarity in Go. If given a tree in newick format, and two text files contain taxa in 
communities A and B (each on one line), Pismo returns the total branch length of community A, total branch length of 
community B, and the total branch length shared by community A and B. From this its trivial to calculate PhyloSor as:

```math
PhylSor = \frac{BL_{shared}}{(BL_{a} + BL_{b}) * 0.5}
```
