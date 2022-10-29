#!/bin/bash
version="v0.1.0";
os=$(uname);
if [ "$os" = "Darwin" ]; then
    os="darwin"
else
    os="linux"
fi

arch=$(uname -m);

if [ "$arch" = "x86_64" ]
then
    arch="amd64"
fi

mkdir cmdr_install;
echo "Installing for cmdr-$version-$os-$arch.tar.gz";
wget https://github.com/alfiankan/commander/releases/download/$version/cmdr-$version-$os-$arch.tar.gz;
tar -xf cmdr-$version-$os-$arch.tar.gz -C cmdr_install;
rm -rf cmdr-$version-$os-$arch.tar.gz;
chmod +x cmdr_install/cmdr;
cp cmdr_install/cmdr /usr/local/bin;
rm -rf cmdr_install;
