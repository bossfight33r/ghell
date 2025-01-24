#!/usr/bin/env python3
import json, subprocess, time, sys, os

WIDTH  = 80
HEIGHT = 24
CWD    = "~/Desktop/pr/ghell"

events = []
t = 0.0

def out(text, delay=0.0):
    global t
    t += delay
    events.append([round(t, 4), "o", text])

def prompt():
    out(f"\r\033[0;32m{CWD}\033[0m \033[1;37m$\033[0m ", 0.05)

def type_cmd(cmd, char_delay=0.07):
    global t
    for ch in cmd:
        t += char_delay
        events.append([round(t, 4), "o", ch])

def enter():
    global t
    t += 0.05
    events.append([round(t, 4), "o", "\r\n"])

def run(cmd, output, pause_before=0.3, pause_after=0.2):
    prompt()
    type_cmd(cmd)
    enter()
    out(output, pause_before)
    out("\r\n", pause_after)


# ── intro pause ──
t += 0.8

# 1. pipes
run(
    "ls | grep -E 'go|shell'",
    "\033[0;36mgo.mod\033[0m\r\ngo.sum\r\n\033[0;36mshell/\033[0m",
)

# 2. quotes
run(
    'mkdir "my folder" && ls | grep my',
    "\033[0;34mmy folder\033[0m",
)
run(
    'rm -r "my folder"',
    "",
    pause_before=0.15,
)

# 3. tilde anywhere
run(
    "ls ~/Desktop | head -3",
    "dbguide.pdf\nfirst-step\nghell",
)

# 4. env vars + exit code
run(
    "NAME=ghell",
    "",
    pause_before=0.1,
)
run(
    'echo "hello from $NAME"',
    "hello from ghell",
)

# 5. && and ||
run(
    "go build ./... && echo build ok",
    "build ok",
    pause_before=0.9,
)
run(
    "cd /bad || echo 'no such dir'",
    "no such dir",
    pause_before=0.2,
)

# 6. background job
run(
    "sleep 2 &",
    "\033[1m[1] 49217\033[0m",
    pause_before=0.1,
)
run(
    "echo doing other things...",
    "doing other things...",
)
out("\r\n[1]+  Done\tsleep 2\r\n", 2.1)

# 7. exit
prompt()
type_cmd("exit")
enter()
t += 0.3

header = {
    "version": 2,
    "width":   WIDTH,
    "height":  HEIGHT,
    "timestamp": int(time.time()),
    "title": "ghell demo",
}

out_path = os.path.join(os.path.dirname(__file__), "demo.cast")
with open(out_path, "w") as f:
    f.write(json.dumps(header) + "\n")
    for ev in events:
        f.write(json.dumps(ev) + "\n")

print(f"wrote {out_path} ({len(events)} events, {t:.1f}s)")
