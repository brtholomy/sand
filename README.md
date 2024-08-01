# sand

Go implementation of the original sandpile experiment showing Self-Organized Criticality (SOC).

To run the simulation with a grid of 50x50 and 1M iterations:

```
go run sand.go --size 50 --iters 1_000_000 --chart
```

![chart](size-50_iters-1M.png)

see all options:

```
go run sand.go --help
```
