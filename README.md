# pokedexcli
Pokedexcli is a command line utility that allows you to query PokeAPI.co for
information about Pokemon. It provides a REPL environment and a number of
commands to get information about game locations and Pokemon. Resources are
cached locally to provide a snappy user experience.

# Installation
First, make sure you have `go` installed (https://go.dev/doc/install). Then,
run `go get https://github.com/caleb-fringer/pokedexcli`. 

# Usage
Start the REPL with `go run .`
Pokedexcli provides several commands for interacting with the API. They are
discoverable with the `help` command. Use `map` and `mapb` to explore locations
available. Use `explore location-name` to get information about the Pokemon in
that location! You can attempt to catch it with `catch pokemon-name`. Once a 
Pokemon has been caught, it may be inspected with `inspect pokemon-name`.

# Demo
<video src="https://github.com/caleb-fringer/pokedexcli/demo.mp4" controls></video>
