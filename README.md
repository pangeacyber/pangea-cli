# Pangea Secrets CLI

The easiest way to scrap .env files and store your API keys securely on [Pangea](https://pangea.cloud?utm_source=github&utm_medium=pangea-cli-repo).

## Installation

### Binary

It is possible to install Pangea CLI directly in binary mode. In order to get it downloaded and installed follow these next steps:


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


### Source code

If prefer to install directly from source code run following commands:

- Clone Pangea CLI repo and move inside its folder
```bash
git clone git@gitlab.com:pangeacyber/pangea-cli.git
cd pangea-cli
```

- Run make installer (only works on macOS and linux)
```bash
make install-<OS>
```

Note: Replace `<OS>` with `macos` or `linux`


## Usage

### Login to Pangea
```bash
pangea login
```
Note: Follow the prompt and paste your Pangea vault token

### Create Workspace
```bash
pangea vault workspace create
```

### Select Workspace
```bash
pangea vault workspace select
```

### Migrate .env file to a Pangea Workspace
```bash
pangea vault workspace migrate -f .env
```

### Run with secrets from Pangea
```bash
pangea vault workspace run -c <APP_COMMAND>
# Example - pangea vault workspace run -c npm run dev
```

### Docker Container

Step 1: Install the CLI in your `Dockerfile`. Here's an example for a Node app
```dockerfile
FROM node:lts-bullseye

# Install Pangea CLI
RUN sh -c "$(curl -fsSL -o pangea.tar.gz https://github.com/pangeacyber/pangea-cli/releases/download/latest/pangea-$(uname -s)-$(uname -m).tar.gz)"
RUN sh -c "mkdir installer && tar -xzvf pangea.tar.gz -C installer && cd installer"
RUN sh -c "./install.sh && cd ../"

WORKDIR /app
COPY . .

RUN npm install

ENTRYPOINT ["pangea", "vault", "workspace", "run", "-c"]
# APP Command
CMD ["npm", "run", "dev"]
```

Now run your docker container by passing in the PANGEA_TOKEN and PANGEA_DEFAULT_FOLDER.
```bash
docker run \
    -e PANGEA_TOKEN=pts... \
    -e PANGEA_DOMAIN=aws.us.pangea.cloud \
    -e PANGEA_DEFAULT_FOLDER=/... \
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


## Develop

Follow this next instructions to test installation script in multiple platforms:

### Linux ARM64

Run this command to build a Linux ARM64 image and run installation script on it:

```bash 
docker build -t pangea-test-linux-arm64 --platform linux/arm64 .  
```

Run docker image and check its logs, it should show the message:

```
Running Pangea-CLI install script...
OS: Linux
Arch: aarch64
go build -o pangea ./cmd/main.go
Bash completion script installed successfully.

Installation success.
```

### Linux AMD64

Run this command to build a Linux AMD64 image and run installation script on it:

```bash 
docker build -t pangea-test-linux-amd64 --platform linux/amd64 .  
```

Run docker image and check its logs, it should show the message:

```
Running Pangea-CLI install script...
OS: Linux
Arch: x86_64
go build -o pangea ./cmd/main.go
Bash completion script installed successfully.

Installation success.
```
