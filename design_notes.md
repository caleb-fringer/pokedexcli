# Design Notes
The current structure is unsustainable as I add more commands to the repl. I
would like to fix this by abstracting away the logic of fetching a resource from
the API into a function that does not care about the type of response back, and
create custom handlers/data structures for parsing each *type*.

One thing I don't like is that the REPL requires every command handler to use
the same set of parameters (PageLink struct). Only the map/mapb seem to make use
of this param and it is otherwise wasted. However, I'm not sure how to specify a
generic function handler. Maybe I can find out more in The Go Programming
Language. Perhaps I could use closures to capture the variables to pass to the
pokeapi functions.

Furthermore, it would seem that the API calling function would be best suited as
a generic function that can take a struct of whatever response type is required
and return that object to the endpoint-specific handler. Or, I could have the
API calling code simply return a byte array and let each individual endpoint
handler worry about creating a struct to hold the object.

I believe only the map/mapb commands need to update the pagelink struct, because
this is state owned by the REPL. Perhaps this state should be owned by the
pokeapi package.
