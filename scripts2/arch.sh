#!/bin/sh

test_arch() {
   local arch
   machine_arch=$1
   case "${machine_arch}" in
       arm*)
           arch=${machine_arch%l}
           ;;
   esac
   printf "machine_arch=%-10s -> arch=%s\n" "$machine_arch" "$arch"
}

test_arch "armv6l"
test_arch "armv6"
test_arch "armv7l"
test_arch "armv7"
test_arch "arm"
test_arch "arml"
test_arch "x86_64"
