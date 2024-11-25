#GO-LLM

basic project, utilizes the go-llama.cpp bindings.


#Usage:
- go-LLM has two modes of interaction, an interactive terminal converstation, or via tcp with JSON requests.
- the listen server is the default mode, the address and port of which can be modified in the config.

- Flags
    - -c this is the path to the config.json, the default is './config.json'
    - --interactive-mode this will start an interactive terminal session.


#TODO:
- ~~TODO: now that it can query a model, a means of persistent conversation with the model.~~
- TODO: front end needed
- ~~TODO: API potentially, to handle queries in a server/client format~~
- TODO: Web scraping - this thing needs the capacity to go poke the web for stuff
- TODO: Research methods of improving context clues, potentially having a _rolling_ history
- TODO: Need to create a better init prompt, but overall function seems to work
- ~~TODO: Add tuning parameters to config, so they are no longer hard-coded.~~
