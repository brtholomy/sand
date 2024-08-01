# sand

Go implementation of the original sandpile experiment showing Self-Organized Criticality (SOC).

To run the simulation with a grid of 50x50 and 1M iterations:

```sh
go run sand.go --size 50 --iters 1_000_000 --chart
```

![chart](chart_50_size_1M_iters.png)

End result of a 500 width pile after 100k iterations:

```sh
go run sand.go -s 500 -i 100_000 --pixel
```

![pile](pile_500px_100k_iters.png)

End result of a 100 width pile after 2M iterations with height of 16:

```sh
go run sand.go -x 16 -s 100 -i 2_000_000 --pixel
```

![pile](pile_100px_2M_iters.png)

see all options:

```sh
go run sand.go --help
```
