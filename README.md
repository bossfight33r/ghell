# ghell

Unix shell written in Go from scratch — no third-party libraries except readline.

![demo](demo.gif)

## What's inside

| Feature | Example |
|---|---|
| Pipelines | `ls \| grep go \| wc -l` |
| I/O redirection | `echo hi > out.txt`, `cat < in.txt` |
| Append | `echo more >> out.txt` |
| Quote handling | `mkdir "my folder"`, `echo 'no $expand'` |
| Tilde expansion | `ls ~/Desktop`, `cd ~/projects` |
| Environment variables | `NAME=Ivan`, `echo $NAME`, `export PATH=$PATH:/bin` |
| Exit codes | `ls /bad; echo $?` |
| Conditional execution | `make build && ./ghell`, `cd /bad \|\| echo nope` |
| Semicolons | `echo one; echo two; echo three` |
| Background jobs | `sleep 5 &`, `jobs` |
| Signal handling | `Ctrl+C` kills the running command, not the shell |
| Tab completion | commands from `$PATH` + files, with `~` support |
| Arrow key history | `↑` / `↓` to navigate, `←` / `→` to edit |

Built-ins: `cd`, `pwd`, `export`, `history`, `jobs`, `exit`

## Architecture

```
ghell/
├── main.go              # REPL loop, readline setup
└── shell/
    ├── shell.go         # Shell struct, Execute, runStep
    ├── tokenize.go      # quote-aware tokenizer, splitAtOps, splitPipeline
    ├── parse.go         # command struct, redirect parsing
    ├── expand.go        # $VAR / $? / ~ expansion, isAssignment
    ├── exec.go          # run single command, I/O redirection wiring
    ├── pipeline.go      # connect N commands via os.Pipe
    ├── builtins.go      # cd, pwd, export, history, jobs, exit
    ├── jobs.go          # background job store (mutex-protected)
    ├── signals.go       # ignore SIGINT/SIGQUIT in shell process
    └── complete.go      # readline tab completer
```

The execution flow for `ls | grep go > out.txt`:

```
Execute()
  └── splitAtOps()          # no &&/||/; here
      └── splitPipeline()   # ["ls", "grep go > out.txt"]
          └── runPipeline()
                ├── parseCommand("ls")             → {args: ["ls"]}
                ├── parseCommand("grep go > out.txt") → {args: ["grep","go"], outFile: "out.txt"}
                ├── os.Pipe() between ls and grep
                ├── p.Start() × 2
                └── p.Wait() × 2
```

## Getting started

```bash
git clone https://github.com/ivan/ghell
cd ghell
make run
```

**Requirements:** Go 1.21+

## Commands

```bash
make build    # compile to ./ghell
make run      # build and launch
make test     # run tests
make lint     # go vet
make install  # install to /usr/local/bin
make clean    # remove binary
```

## Tests

```bash
make test
```

Covers: argument parsing, all redirect operators, pipeline splitting, variable expansion (`$VAR`, `$?`, undefined vars, edge cases), `isAssignment` validation.
