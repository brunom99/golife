A simple Go project to explore websockets, goroutines and mutex.

Each goroutine is represented by a bubble moving across the grid.

There are 3 types of bubble: common, light and dark.

- Dark bubbles destroy common bubbles and replicate themselves.
- Light bubbles destroy dark bubbles.

Execute the command "go run main.go" and navigate to http://localhost:1984/

You can modify parameters in the "config.toml" configuration file.

<img width="861" alt="01" src="https://github.com/user-attachments/assets/2164da64-ec49-4253-855f-32240d299d27">
