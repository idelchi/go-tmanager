rm -rf ~/.local/bin
mkdir -p ~/.local/bin

rm -rf ~/.bin-*
rm -rf tests
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/scripts/install.sh | sh -s -- -v v0.0-dev -o ~/.local/bin

~/.local/bin/godyl --dot-env="${HOME}"/.secrets/github/.github --tags pi5-64,pi "${HOME}"/.pi/tools/tools.yml --output=~/.local/bin

rm -rf ~/.local/bin
mkdir -p ~/.local/bin

rm -rf ~/.bin-*
rm -rf tests
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/scripts/install.sh | sh -s -- -v v0.0-dev -o ~/.local/bin
~/.local/bin/godyl --dot-env="${HOME}"/.secrets/github/.github --tags pi0-64,pi "${HOME}"/.pi/tools/tools.yml --output=~/.local/bin

rm -rf ~/.local/bin
mkdir -p ~/.local/bin

rm -rf ~/.bin-*
rm -rf tests
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/scripts/install.sh | sh -s -- -v v0.0-dev -o ~/.local/bin
~/.local/bin/godyl --dot-env="${HOME}"/.secrets/github/.github --tags pi5-32,pi "${HOME}"/.pi/tools/tools.yml --output=~/.local/bin

cd /tmp || exit

curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/scripts/install.sh | sh -s -- -v v0.1-beta -o .

./godyl --detect
