# Pangea Secrets CLI

The easiest way to scrap .env files and store your API keys securely on [Pangea](https://pangea.cloud?utm_source=github&utm_medium=pangea-cli-repo).

## Get Started
[Video Walkthrough on Getting Started](https://www.youtube.com/watch?v=R_LSoDcXj9Y)

### Installation

#### Using `brew` run:
```bash
brew install pangeacyber/cli/pangea
```

#### For linux / macOS systems

using `curl`
```bash
source <(curl -L https://github.com/pangeacyber/pangea-cli/releases/latest/download/download-and-install.sh)
```

or using `wget`
```bash
bash <(wget -qO- https://github.com/pangeacyber/pangea-cli/releases/latest/download/download-and-install.sh)
```

#### For Windows systems

Using `winget`, run:
```bash
winget install pangeacyber.pangea
```



### Login to Pangea
```bash
pangea login
```
Note: Follow the prompt and paste your Pangea vault token

### Create Workspace
```bash
pangea create
```

### Select Workspace
```bash
pangea select
```

### Migrate .env file to a Pangea Workspace
```bash
pangea migrate -f .env
```

### Run with secrets from Pangea
```bash
pangea run -- <APP_COMMAND>
# Example - pangea run -- npm run dev
```

## Usage
### Docker Container
Step 1: Install the CLI in your `Dockerfile`. Here's an example for a Node app
```dockerfile
FROM node:lts-bullseye

# Install Pangea CLI
RUN curl -L -o /bin/pangea "https://github.com/pangeacyber/pangea-cli/releases/latest/download/pangea-$(uname -s)-$(uname -m)"

WORKDIR /app
COPY . .

RUN npm install

ENTRYPOINT ["pangea", "run", "-c"]
# APP Command
CMD ["npm", "run", "dev"]
```

Now run your docker container by passing in the PANGEA_TOKEN and PANGEA_DEFAULT_FOLDER.
```bash
docker run \
    -e PANGEA_TOKEN=pts... \
    -e PANGEA_DOMAIN=aws.us.pangea.cloud \
    -e PANGEA_DEFAULT_FOLDER=/secrets/... \
    <IMAGE_NAME>
```


## Troubleshooting

### Completion

If completion is not working on MacOS using `zsh`, you could try adding this next lines to your `~/.zshrc` file

```
if type brew &>/dev/null
then
  FPATH="$(brew --prefix)/share/zsh/site-functions:${FPATH}"

  autoload -Uz compinit
  compinit
fi
```

and then restart the terminal running: 

```sh
exec $SHELL
```
