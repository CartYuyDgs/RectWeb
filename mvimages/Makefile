NAME=agent
OS=linux
#OS=windows
# ARCH=amd64
# ARCH=arm
ARCH=arm64
ARM=7



arm64:
        GOOS=$(OS)  GOARCH=$(ARCH)  GOARM=$(ARM)  go build -o $(NAME)_$(OS)_$(ARCH)_$(ARM) cmd/bonc_ai/main.go