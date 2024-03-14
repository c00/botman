# Botman

Botman lets you talk to an LLM. Currently only openAi GPT 4 turbo It is optimized for use in the terminal. Can use `stdin` or an argument for input, outputs content to `stdout` and errors to `stderr`.

## Setup

```
go install github.com/c00/botman
```

Set an environment variable called `OPENAI_API_KEY` with your api key.

## Examples

```bash
# Use an argument as input.
botman "tell me a joke about the golang gopher"

# Use stdin as input
echo Quote a Bob Kelso joke | botman
```
