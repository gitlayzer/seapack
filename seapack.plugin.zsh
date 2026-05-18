# SeaPack Zsh Plugin
# Loads completions dynamically from the seapack binary

if (( $+commands[seapack] )); then
    source <(seapack completion zsh)
fi
