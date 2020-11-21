.PHONY : cli

build : cli

install : cli
		cd cli && make install

cli :
		cd cli && make

clean:
		cd cli && make clean
