
pkgname=rollout-cli
pkgver=tbd
pkgrel=1
pkgdesc="Rollout CLI"
arch=(x86_64)
url="https://github.com/pjeziorowski/rollout"
license=('MIT')
makedepends=(git go)
source=("https://github.com/pjeziorowski/rollout/archive/v$pkgver.tar.gz")

build() {
	cd "$pkgname-$pkgver"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
    go build -o $pkgname main.go
}

package() {
	cd "$pkgname-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/rollout"
}