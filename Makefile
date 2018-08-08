build:
	mkdir -p bin
	cd xfs-watch && go build && cd -
	mv xfs-watch/xfs-watch bin
