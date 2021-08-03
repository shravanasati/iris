
#!/bin/bash

echo "Downloading iris..."
curl -L "https://github.com/Shravan-1908/iris/releases/latest/download/iris-darwin-amd64" -o iris

echo "Adding iris into PATH..."

mkdir -p ~/.iris;
mv ./iris ~/.iris
echo "export PATH=$PATH:~/.iris" >> ~/.bashrc

echo "iris installation is completed!"
echo "You need to restart the shell to use iris."
