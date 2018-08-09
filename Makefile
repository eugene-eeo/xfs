build:
	mkdir -p bin
	cd cmd/xfs-watch && go build && cd -
	cd cmd/xfs-index && go build && cd -
	cd cmd/xfs-search && go build && cd -
	cd cmd/xfs-dispatch && go build && cd -
	mv cmd/xfs-watch/xfs-watch bin
	mv cmd/xfs-index/xfs-index bin
	mv cmd/xfs-search/xfs-search bin
	mv cmd/xfs-dispatch/xfs-dispatch bin
