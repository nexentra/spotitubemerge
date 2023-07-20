#!/bin/bash

# Check if any argument is provided
if [ $# -eq 0 ]; then
    # If no argument is provided, run both frontend and backend
    run_frontend=true
    run_backend=true
else
    # Check if the argument is 'frontend'
    if [ "$1" == "frontend" ]; then
        run_frontend=true
        run_backend=false
    # Check if the argument is 'backend'
    elif [ "$1" == "backend" ]; then
        run_frontend=false
        run_backend=true
    else
        echo "Invalid argument. Usage: $0 [frontend|backend]"
        exit 1
    fi
fi

# Run tmux with the specified configuration
tmux new-session -d bash
tmux split-window -h bash

# Run frontend if specified
if [ "$run_frontend" = true ]; then
    tmux send -t 0:0.0 "make run-frontend" C-m
fi

# Run backend if specified
if [ "$run_backend" = true ]; then
    tmux send -t 0:0.1 "make run-backend" C-m
fi

tmux -2 attach-session -d
