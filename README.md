# AIChat

AIChat is a TUI chat application that utilizes OpenAI text completions. Go ahead, have a conversation.

## Installation

You can download the latest version from the [releases page](https://github.com/KalebHawkins/aichat/releases).

## Usage

You must get an API Key from OpenAI. Once you have that export an `AI_CHAT_KEY` environment variable.

For Linux or Mac users:

```bash
export AI_CHAT_KEY="OPENAI_API_KEY"
```

For Windows users:

```powershell
$env:AI_CHAT_KEY="OPENAI_API_KEY"
```

Now you can run the downloaded binary, type your question and hit <kbd>Ctrl</kbd>+<kbd>S</kbd> to submit your message.

You will get a response fairly quickly depending on the complexity of your input.

The input may be incomplete, this is due to the way the token usage is set up. By default, you may submit a receive a query that uses up to 500 tokens. If you want to complete a response all you need to do is <kbd>Ctrl</kbd>+<kbd>S</kbd> to submit the query with the output. 

If you want to quickly start a fresh session you can clear the text area by pressing <kbd>Ctrl</kbd>+<kbd>C</kbd>. 

To exit you can use the <kbd>Esc</kbd> key.





