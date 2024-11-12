source .env
alias g="godyl --log debug"
alias gr="go run . --log=debug"

echo "done!"
alias d="go run . --dry --log debug"


export $(grep -v '^#' taskfile.env | xargs)
