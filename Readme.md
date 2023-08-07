# Pangea Secrets CLI

The easiest way to scrap .env files and store your API keys securely on [Pangea](https://pangea.cloud).

## Get Started

### Installation
For linux / macOS systems
```bash
curl -L -o /usr/local/bin/pangea "https://github.com/pangeacyber/pangea-cli/releases/latest/download/pangea-$(uname -s)-$(uname -m)"
```

### Login to Pangea
```bash
pangea login
### Follow the prompt and paste your Pangea vault token
```

### Create Project
```bash
pangea create
```

### Select Project
```bash
pangea select
```

### Migrate .env file to a Pangea project
```bash
pangea migrate -f .env
```

### Run with secrets from Pangea
```bash
pangea run -c <APP_COMMAND>
# Example - pangea run -c npm run dev
```

## Usage
### Docker Container
Step 1: Install the CLI in your `Dockerfile`. Here's an example for a Node app
```docker
FROM node:lts-bullseye

# Install Pangea CLI
curl -L -o /usr/local/bin/pangea "https://github.com/pangeacyber/pangea-cli/releases/latest/download/pangea-$(uname -s)-$(uname -m)"

WORKDIR /app
COPY . .

RUN npm install

ENTRYPOINT ["pangea", "run", "-c"]
# APP Command
CMD ["npm", "run", "dev"]
```

Now run your docker container by passing in the PANGEA_TOKEN and PANGEA_DEFAULT_PATH.
```bash
docker run \
    -e PANGEA_TOKEN=pts... \
    -e PANGEA_DEFAULT_PATH=/secrets/... \
    <IMAGE_NAME>
```