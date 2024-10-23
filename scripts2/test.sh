rm -rf .bin-*
ARGS="--dry --log=info"

echo "-------- os: linux, arch: amd64"
go run ./cmd/godyl ${ARGS} --os=linux --arch=amd64 --output=".bin-{{ .OS }}-{{ .ARCH }}"
echo "-------- os: linux, arch: arm64"
go run ./cmd/godyl ${ARGS} --os=linux --arch=arm64 --output=".bin-{{ .OS }}-{{ .ARCH }}"
echo "-------- os: linux, arch: armv7"
go run ./cmd/godyl ${ARGS} --os=linux --arch=armv7 --output=".bin-{{ .OS }}-{{ .ARCH }}{{ .ARCH_VERSION }}"
echo "-------- os: linux, arch: armv6"
go run ./cmd/godyl ${ARGS} --os=linux --arch=armv6 --output=".bin-{{ .OS }}-{{ .ARCH }}{{ .ARCH_VERSION }}"
echo "-------- os: darwin, arch: amd64"
go run ./cmd/godyl ${ARGS} --os=darwin --arch=amd64 --output=".bin-{{ .OS }}-{{ .ARCH }}"
echo "-------- os: darwin, arch: arm64"
go run ./cmd/godyl ${ARGS} --os=darwin --arch=arm64 --output=".bin-{{ .OS }}-{{ .ARCH }}"
echo "-------- os: windows, arch: amd64"
go run ./cmd/godyl ${ARGS} --os=windows --arch=amd64 --output=".bin-{{ .OS }}-{{ .ARCH }}"

