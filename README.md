# A smarter email assistant using LLM

[[`Project Writeup`](https://medium.com/@jonathanlawhh) [`My Website`](https://jonathanlawhh.com/)]

## Project Overview

Email AI auto replies are great today, but they are not customized. Certain people should receive certain way of replying (for example my manager).

What if we build a simple too, that allows you to define rules on how the AI respond to the email.


## References

- [Outlook Co-Pilot](https://create.microsoft.com/en-us/learn/articles/use-copilot-to-draft-email-replies-in-outlook)
- [Gmail Smart Compose](https://blog.google/products/gmail/gmail-ai-features/)
- [Gemini API](https://ai.google.dev/)

## Setup and Usage

### Software Requirements

- Golang
- [Gemini API key](https://ai.google.dev/gemini-api/docs/api-key)

### Installation

1. Clone this repository:

```bash
git clone https://github.com/jonathanlawhh/a-smarter-email-assistant
```

2. Install required libraries:

```bash
go install
```

### Usage

1. Setup the `.env` with required information from `.env-sample`

2. Run the script.

```bash
go run .\Local\main.go
```

`.env` parameters:

| ENV NAME            | Accepted values    | Description                                  |
|---------------------|--------------------|----------------------------------------------|
| GEMINI_API_KEY      | string             | Gemini API Key                               |
              |

## Closing thoughts

- Let me think...
- To add OpenAI endpoint too