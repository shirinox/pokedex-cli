module github.com/shirinox/pokedex

go 1.22.4

replace github.com/shirinox/pokeapi v0.0.0 => ./internal/pokeapi

replace github.com/shirinox/pokecache v0.0.0 => ./internal/pokecache

require (
	github.com/shirinox/pokeapi v0.0.0
	github.com/shirinox/pokecache v0.0.0
)
