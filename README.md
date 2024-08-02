# sand

Go implementation of the original sandpile experiment showing Self-Organized Criticality (SOC).

To run the simulation with a grid of 50x50 and 1M iterations:

```sh
go run sand.go --size 50 --iters 1_000_000 --chart
```

![chart](imgs/chart_50_size_1M_iters.png)

End result of a 500 width pile after 100k iterations:

```sh
go run sand.go -s 500 -i 100_000 --pixel
```

![pile](imgs/pile_500px_100k_iters.png)

End result of a 200 width pile after 4M iterations with height of 12:

```sh
go run sand.go -x 12 -s 200 -i 4_000_000 --pixel
```

![pile](imgs/pile_200-size_4000000-iters_12-height.png)

see all options:

```sh
go run sand.go --help
```
