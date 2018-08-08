build:
	mkdir -p bin
	cd cmd/xfs-watch && go build && cd -
	cd cmd/xfs-index && go build && cd -
	mv cmd/xfs-watch/xfs-watch bin
	mv cmd/xfs-index/xfs-index bin
