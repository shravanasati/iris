#!/usr/bin/bash

platforms=("windows/amd64" "darwin/amd64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="./bin/iris-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "==> Building executable for $GOOS/$GOARCH..."

    env GOOS=$GOOS GOARCH=$GOARCH go build -v -ldflags="-s -w" -o $output_name ./cmd/iris
    if [ $? -ne 0 ]; then
        echo '==> An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "==> Executables built successfully!"