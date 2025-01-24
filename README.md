# ghell

A Unix shell written in Go from scratch.

![demo](demo.gif)

## Features

- **Pipelines** — `ls | grep go | wc -l`
- **I/O redirection** — `>`, `>>`, `<`
- **Environment variables** — `NAME=Ivan`, `echo $NAME`, `export`
- **Quote handling** — `mkdir "my folder"`, `echo 'literal $VAR'`
- **Tilde expansion** — `ls ~/Desktop` works anywhere
- **Operators** — `&&`, `||`, `;`
- **Background jobs** — `sleep 5 &`, `jobs`
- **Exit codes** — `echo $?`
- **Signal handling** — `Ctrl+C` kills the current command, not the shell
- **readline** — arrow key history, line editing, tab completion
- **Built-ins** — `cd`, `pwd`, `history`, `export`, `jobs`, `exit`

## Run

```bash
make run
```

## Test

```bash
make test
```

## Install

```bash
make install   # puts ghell in /usr/local/bin
```
