#!/bin/sh
set -eu

SESSION="aniways"
ROOT="${PWD}"

if ! command -v tmux >/dev/null 2>&1; then
  echo "tmux is not installed. Please install tmux and try again." >&2
  exit 1
fi

if tmux has-session -t "$SESSION" 2>/dev/null; then
  if [ -n "${TMUX-}" ]; then
    exec tmux switch-client -t "$SESSION"
  else
    exec tmux attach -t "$SESSION"
  fi
fi

tmux new-session -d -s "$SESSION" -n editor -c "$ROOT"
tmux send-keys -t "$SESSION:0" 'nvim .' C-m

tmux new-window  -t "$SESSION:1" -n term   -c "$ROOT"

tmux new-window  -t "$SESSION:2" -n api    -c "$ROOT"
tmux send-keys   -t "$SESSION:2" 'make dev-api' C-m

tmux new-window  -t "$SESSION:3" -n worker -c "$ROOT"
tmux send-keys   -t "$SESSION:3" 'make dev-worker' C-m

tmux new-window  -t "$SESSION:4" -n proxy  -c "$ROOT"
tmux send-keys   -t "$SESSION:4" 'make dev-proxy' C-m

tmux new-window  -t "$SESSION:5" -n web    -c "$ROOT"
tmux send-keys   -t "$SESSION:5" 'cd web && bun run dev' C-m

tmux select-window -t "$SESSION:0"
if [ -n "${TMUX-}" ]; then
  exec tmux switch-client -t "$SESSION"
else
  exec tmux attach -t "$SESSION"
fi
