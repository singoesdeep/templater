#!/bin/bash

set -e

# Create necessary directories
mkdir -p packages/{deb,rpm,apk}

# Create Debian package structure
cat > packages/deb/control << EOF
Package: templater
Version: 0.1.0
Section: utils
Priority: optional
Architecture: amd64
Depends: libc6
Maintainer: Your Name <your.email@example.com>
Description: A template processing tool
 templater is a Go-based template processing tool designed to generate
 code and other text-based files from templates using data from various sources.
EOF

# Create RPM package structure
cat > packages/rpm/templater.spec << EOF
Name:           templater
Version:        0.1.0
Release:        1%{?dist}
Summary:        A template processing tool

License:        MIT
URL:            https://github.com/singoesdeep/templater
BuildArch:      x86_64

%description
templater is a Go-based template processing tool designed to generate
code and other text-based files from templates using data from various sources.

%prep
# Nothing to do here

%build
# Nothing to do here

%install
mkdir -p %{buildroot}/usr/bin
cp templater %{buildroot}/usr/bin/

%files
/usr/bin/templater

%changelog
* $(date '+%a %b %d %Y') Your Name <your.email@example.com> - 0.1.0-1
- Initial release
EOF

# Create Alpine package structure
cat > packages/apk/APKBUILD << EOF
# Contributor: Your Name <your.email@example.com>
# Maintainer: Your Name <your.email@example.com>
pkgname=templater
pkgver=0.1.0
pkgrel=0
pkgdesc="A template processing tool"
url="https://github.com/singoesdeep/templater"
arch="noarch"
license="MIT"
depends=""
makedepends=""
subpackages=""
options=""
builddir="\$srcdir/\$pkgname-\$pkgver"

package() {
    mkdir -p "\$pkgdir"/usr/bin
    cp templater "\$pkgdir"/usr/bin/
}
EOF

echo "Package repository structure created in packages/"
echo "Don't forget to update the maintainer information in the package files." 