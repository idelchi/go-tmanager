docker compose run --build bookworm-amd64 godyl --detect
docker compose run --build bookworm-amd32 godyl --detect
docker compose run --build bookworm-arm64 godyl --detect
docker compose run --build bookworm-armv7 godyl --detect
docker compose run --build bookworm-armv5 godyl --detect

docker compose run --build alpine-arm64 godyl --detect
docker compose run --build alpine-armv7 godyl --detect
docker compose run --build alpine-armv6 godyl --detect
docker compose run --build alpine-armv5 godyl --detect
go run ./cmd/godyl --detect


rm -rf .bin-* .logs
mkdir -p .logs

# Debian Bookworm builds
docker compose run --build bookworm-amd64 bash -c "godyl --log=info > .logs/bookworm-amd64.log 2>&1"
docker compose run --build bookworm-arm64 bash -c "godyl --log=info > .logs/bookworm-arm64.log 2>&1"
docker compose run --build bookworm-armv7 bash -c "godyl --log=info > .logs/bookworm-armv7.log 2>&1"
docker compose run --build bookworm-armv5 bash -c "godyl --log=info > .logs/bookworm-armv5.log 2>&1"

# Alpine builds
docker compose run --build alpine-arm64 ash -c "godyl --log=info > .logs/alpine-arm64.log 2>&1"
docker compose run --build alpine-armv7 ash -c "godyl --log=info > .logs/alpine-armv7.log 2>&1"
docker compose run --build alpine-armv6 ash -c "godyl --log=info > .logs/alpine-armv6.log 2>&1"
go install ./cmd/godyl
godyl --log=info --arch=armv5 > .logs/windows-amd64.log 2>&1



