# Botman

Botman lets you talk to an LLM. Currently only openAi GPT 4 turbo It is optimized for use in the terminal. Can use `stdin` or an argument for input, outputs content to `stdout` and errors to `stderr`.

## Setup

1. Clone the repo
2. run `go install .`
3. Add environment variable `OPENAI_API_KEY` with your api key.
4. Test that it works by running `botman -h`

## Examples

```bash
# Use an argument as input.
botman "tell me a joke about the golang gopher"
botman "git command for undoing last commit"
botman "untar .tar.gz file into new directory"

# Use stdin as input
echo Quote a Bob Kelso joke | botman

# Use both
ls -al | botman "Which files are hidden?"
cat deployment.yaml | botman "how many replicas will this run?"
```

![demo](https://github.com/c00/botman/blob/main/assets/botman-demo.gif?raw=true)

## Interactive mode

Use interactive mode by giving the `-i` flag.

In interactive mode, the program does not exit after a response, but instead waits for your input.

```bash
# Start interactive mode without an initial prompt
bot -i

# Start interactive mode with arguments
bot -i How many bees in a bonnet?
```
