## Prompt AI

Simple TUI prompt aiming at integrating with Chat GPT and more.
![TUI example](https://github.com/noboruma/prompt-ai/blob/main/.github/presentation.png?raw=true)

## Prerequisites

To use this tool, you will need to provide your Open AI API Key.
You can generate one [here](https://platform.openai.com/account/api-keys).

You will need to setup a payment method to make it usable.
This should not cost anything, as long as you are not requesting Chat GPT too much.

## How to install

```
$ go install github.com/noboruma/prompt-ai@latest
```

## Usage
Find your key `<API_KEY>` [here](https://platform.openai.com/account/api-keys)
Then run:

```
$ OPENAI_API_KEY=<API_KEY> prompt-ai
```
or
```
$ export OPENAI_API_KEY=<API_KEY>
$ prompt-ai
```
