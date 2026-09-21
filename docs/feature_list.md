We will aim for concurrency where it makes sense this will not be explicit in the checklist we will discuss this as we start implementing.  

- [ ] Language parsing
    - [ ] Multi-language parsing (C#)
    - [ ] AST based parser (most likely)
    - [ ] Better than regex parsing? (will have to research)
    - [ ] go-tree-sitter which gives us access to tree-sitter this should solve most of the parsing problems

- [ ] CLI Entry (commands)
    - [ ] init
    - [ ] render
    - [ ] render-diff
    - [ ] create-action

- [ ] Change detection (diffing) (C#)
    - [ ] Caching system to allow for partial updates of graphs
    - [ ] Gitsnapshots?
    - [ ] Reimplement the change detector system from C# 

- [ ] Graph modeling
    - [ ] BT graph in python (we will have to look at Go dependencies for this)
    - [ ] Merge/overwrite behavior that interacts with the parsers (C#)

- [ ] Rendering engine
    - Currently uses plantuml but we might switch to mermaid

- [ ] Utils
    - [x] Path managers (more than likely will change to glob)
    - [x] Input handling (more than likely change to glob based structure)
    - [ ] Serealizers 

- [ ] Config managment
    - [ ] Viper based in Go


