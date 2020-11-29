BUILDID=$(git rev-parse --short HEAD)
BUILDOS=linux
GOOS=$BUILDOS go build -o pigwar.server.$BUILDID.$BUILDOS
if [ $? -ne 0 ]; then
   echo "BUILD FAILED"
   exit 1
fi
