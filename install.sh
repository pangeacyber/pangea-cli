#!/bin/bash
BINARY_NAME="pangea"

echo "Running Pangea-CLI install script..."
echo "OS: $(uname -s)"
echo "Arch: $(uname -m)"

case "$SHELL" in
  *zsh*)
    BASH_FILE=".zshrc"
    SOURCE_COMMAND="autoload -U compinit; compinit"
    ;;
  *bash*)
    BASH_FILE=".bashrc"
    SOURCE_COMMAND="source /etc/bash_completion.d/pangea-cli.bash"
    ;;
  *)
    echo "not supported shell"
    exit 1
  ;;
esac

BASH_PATH="${HOME}/${BASH_FILE}"
INSTALL_PATH=/usr/local/bin/

# Check if the file named "pangea" exists in the specified directory
if [ ! -f "./${BINARY_NAME}" ]; then
  echo "There is no binary called '${BINARY_NAME}' in this folder."
  exit 1
fi

grep -Fxq "${SOURCE_COMMAND}" ${BASH_PATH} || echo "${SOURCE_COMMAND}" >> ${BASH_PATH}

# Copy to bin folder
sudo cp ${BINARY_NAME} ${INSTALL_PATH}

# Make it executable
sudo chmod +x ${INSTALL_PATH}${BINARY_NAME}

setup_zsh_completion() {
  if ! ${INSTALL_PATH}${BINARY_NAME} completion zsh > /tmp/${BINARY_NAME}.zsh; then
    echo "Failed to generate completion script"
    return 1
  fi

  # Get the first element of fpath using Zsh
  autoload_dir=$(zsh -c 'echo $fpath[1]')
  sudo mkdir -p ${autoload_dir}

  if ! sudo cp /tmp/${BINARY_NAME}.zsh "${autoload_dir}/_${BINARY_NAME}.zsh"; then
    echo "Failed to move autocompletion script"
    return 1
  fi

  echo "Zsh completion script installed successfully"
}

# Function to handle Bash completion setup
setup_bash_completion() {
  if ! ${INSTALL_PATH}${BINARY_NAME} completion bash > /tmp/${BINARY_NAME}.bash; then
    echo "Failed to generate completion script"
    return 1
  fi

  if ! sudo cp /tmp/${BINARY_NAME}.bash /etc/bash_completion.d/${BINARY_NAME}.bash; then
    echo "Failed to move completion script"
    return 1
  fi

  echo "Bash completion script installed successfully."
}

# completion setup
case "$SHELL" in
  *zsh*)
    setup_zsh_completion
    ;;
  *bash*)
    setup_bash_completion
    ;;
  *)
    echo "Unsupported shell."
    exit 1
    ;;
esac

echo ""
echo "Installation success."
