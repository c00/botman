# botman

`botman` lets you talk to an LLM. Currently it uses openAI ChatGPT-4o. It is optimized for use in the terminal. You can use `stdin` or an argument for input, outputs content to `stdout` and errors to `stderr`.

Note that `botman` does not run any commands. It does not have the ability to _do_ anything, as having some automated LLM control your terminal couldn't possibly lead to anything good. So `botman` simply shows you the output.

## Install from source

1. Clone the repo
2. Run `go install .`
3. Add environment variable `OPENAI_API_KEY` with your api key. e.g. `echo 'export OPENAI_API_KEY="yourapikey"' >> ~/.bashrc`
4. (optional) Create an alias in your shell. e.g. `echo 'alias bot="botman"' >> ~/.bashrc`
5. Test that it works by running `botman "say hi"` or `bot "say hi"`

## Examples

```bash
# Open a new interactive chat.
botman

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

In interactive mode, the program does not exit after a response, so you can continue the conversation.

By not supplying any arguments or stdin, `botman` will run in interactive mode.

Use interactive mode explicitly combined with stdin or arguments by giving the `-i` flag.

```bash
# Start interactive mode without an initial prompt
botman

# Start interactive mode with arguments
botman -i "How many bees in a bonnet?"
```

## Data privacy

`botman` talks directly to the OpenAi API. So assume that OpenAi knows about your plans to overthrow goverments and such. Other than that, botman does not reach out to any service. It currently does not store any information locally either. Tho there are plans to keep a local history for convenience.

## Motivation

I created it mainly for myself but thought it might be useful for others. My motivation stems from seeing some closed-source CLI-LLM integrations from companies I don't necessarily trust. So, I created something free and open source for those of us who value open source and transparency. (Yes, it does still use OpenAI's API, but I am working towards abstracting that away so it could use any LLM.)

# Roadmap

I'm adding features as I feel they're useful. I'm open to suggestions and outside contributions. The aim is to be simple, non-intrusive and transparent about data.

- [ ] History - Store conversation history locally in text files so users can continue older conversations and replay earlier responses.
- [ ] LLM agnostic - Make botman able to work with any LLM by abstracting the interface to the LLM.
- [ ] Ability to execute or at least copy to clipboard the last printed command.
