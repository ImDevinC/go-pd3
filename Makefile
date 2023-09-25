build-linux-x64:
	mkdir -p out/linux-x64
	GOOS=linux GOARCH=amd64 go build -o out/linux-x64/pd3-challenges ./main.go

build-osx-x64:
	mkdir -p out/osx-x64
	GOOS=darwin GOARCH=amd64 go build -o out/osx-x64/pd3-challenges ./main.go

build-windows-x64:
	mkdir -p out/windows-x64
	GOOS=windows GOARCH=amd64 go build -o out/windows-x64/pd3-challenges.exe ./main.go

release: build-linux-x64 build-osx-x64
	tar czvf out/pd3-challenges-${RELEASE_VERSION}-linux-x64.tar.gz --directory out/linux-x64/ pd3-challenges
	tar czvf out/pd3-challenges-${RELEASE_VERSION}-osx-x64.tar.gz --directory out/osx-x64/ pd3-challenges
	zip out/pd3-challenges-${RELEASE_VERSION}-windows-x64.tar.gz --directory out/windows-x64/ pd3-challenges

clean:
	rm -rf out/