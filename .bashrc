source .env
alias g="godyl --log debug"
alias gr="go run ./cmd/godyl --log=debug"

echo "done!"
alias d="go run ./cmd/godyl --dry --log debug"


export $(grep -v '^#' taskfile.env | xargs)
