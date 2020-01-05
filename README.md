# distributable-todo-go

A sample implementation of the [Raft algorithm](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=2ahUKEwis0JeTp6bmAhWjiOAKHXIRBfwQFjAAegQIAxAG&url=https%3A%2F%2Fraft.github.io%2F&usg=AOvVaw0gbPkPuRwWu0Kd74PJmOzK).

# Build

Run `go build ./...`

# Update

Run `go get -u ./...`

# Test

Run `DCONFIG_DIR=$(pwd) go test -v --tags integration ./...`
