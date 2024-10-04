BINARY_NAME=pangea
MAIN=./cmd/main.go
INSTALL_FOLDER=/usr/local/bin/

all: clean build

install-macos: clean build _install-msg-start _install-zsh-cmd _msg-done-restart

uninstall-macos: _uninstall-msg-start completion-zsh-uninstall _uninstall-cmd-zsh _uninstall-msg-done

install-linux: clean build _install-msg-start _install-linux-cmd _msg-done-restart

uninstall-linux: _uninstall-msg-start completion-linux-uninstall _uninstall-cmd-linux _uninstall-msg-done

build-all:
	mkdir -p bin/
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-x86_64 ${MAIN}
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-x86_64 ${MAIN}
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-x86_64.exe ${MAIN}
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-arm64 ${MAIN}
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-arm64 ${MAIN}
	GOARCH=arm64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-arm.exe ${MAIN}

build:
	go build -o ${BINARY_NAME} ${MAIN}

run: build
	./${BINARY_NAME}

verify:
	go vet --all ./...

completion-linux: _check-root
	# Generate completion script
	go run ${MAIN} completion bash > /tmp/${BINARY_NAME}.bash

	# Move to the autoload folder (usually /etc/bash_completion.d for system-wide or ~/.bash_completion for user)
	@sudo cp /tmp/${BINARY_NAME}.bash /etc/bash_completion.d/${BINARY_NAME}

	# # Enabled autoload if not in .bashrc
	@grep -Fxq "source /etc/bash_completion.d/${BINARY_NAME}" $(HOME)/.bashrc || echo "source /etc/bash_completion.d/${BINARY_NAME}" >> $(HOME)/.bashrc

	@echo 'Done. Restart shell to apply changes running "exec $$SHELL"'

completion-linux-uninstall: _check-root
	@echo 'Uninstalling completion on linux'
	@sudo rm -f /etc/bash_completion.d/${BINARY_NAME}
	@echo 'Done. Restart shell to apply changes running "exec $$SHELL"'

completion-zsh:
	# Enabled autoload if not in .zshrc
	@grep -Fxq "autoload -U compinit; compinit" $(HOME)/.zshrc || echo "autoload -U compinit; compinit" >> $(HOME)/.zshrc

	# Generate completion script
	go run ${MAIN} completion zsh > /tmp/${BINARY_NAME}.zsh

	# Move to the autoload folder
	@zsh -c 'sudo cp /tmp/${BINARY_NAME}.zsh $${fpath[1]}/_${BINARY_NAME}'
	@echo 'Done. Restart shell to apply changes running "exec $$SHELL"'

completion-zsh-uninstall:
	@echo 'Uninstalling completion zsh'
	@zsh -c 'sudo rm -f $${fpath[1]}/_${BINARY_NAME}'
	@echo 'Done. Restart shell to apply changes running "exec $$SHELL"'

dev:
	go run ${MAIN}

clean:
	rm -rf bin/*
	rm -rf ~/.pangea/cache
	rm -f ./${BINARY_NAME}
	rm -f ./cmd/${BINARY_NAME}

integration: build
	go test -count=1 -v ./cmd

# Hidden
_install-msg-start:
	@echo "Installing ${BINARY_NAME} command..."

_msg-done-restart:
	@echo 'Done. Restart shell to apply changes running "exec $$SHELL"'

_install-zsh-cmd:
	@zsh -c 'sudo cp ${BINARY_NAME} ${INSTALL_FOLDER}${BINARY_NAME}'

_install-linux-cmd:
	sudo cp ${BINARY_NAME} ${INSTALL_FOLDER}${BINARY_NAME}

_uninstall-msg-start:
	@echo "Uninstalling Pangea CLI."

_uninstall-msg-done:
	@echo "Done. Pangea CLI uninstalled."

_uninstall-cmd-zsh:
	@zsh -c 'sudo rm -f ${INSTALL_FOLDER}${BINARY_NAME}'

_uninstall-cmd-linux:
	sudo rm -f ${INSTALL_FOLDER}${BINARY_NAME}

_check-root:
	@if [ "$$EUID" -ne 0 ]; then \
		echo "Please run as root or using sudo"; \
		exit 1; \
	fi
