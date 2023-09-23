# The ip:port or domain.tld for the server. Don't include protocols.
BASE_LOCATION=""
# Keep this as 's', unless you're testing on a server without ssl
SECURE_S="s"

go build -ldflags "-s -w -X 'main.AUTHOR_NAME=$AUTHOR_NAME' -X 'main.AUTHOR_COLOR=$AUTHOR_COLOR' -X 'main.BASE_LOCATION=$BASE_LOCATION' -X 'main.SECURE_S=$SECURE_S'"