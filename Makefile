# This is the build script.
# It provides commands for building all three projects, together or separately, storing them in 'bin' as independent binary executables.
# Running 'make' compiles everything, and 'minvm' in both debug and release modes.

# Usage:

# make                - builds everything
# make clean          - removes all built binaries
# make min            - just builds 'min'
# make minc           - just builds 'minc'
# make minvm          - builds 'minc' in both debug and release mode
# make minvm-debug    - builds 'minvm' in debug mode
# make minvm-release  - builds 'minvm' in release mode

OUT := bin

all: $(OUT)/min $(OUT)/minc $(OUT)/minvm

$(OUT):
	mkdir -p $(OUT)

$(OUT)/min: | $(OUT)
	cd min && go build -o ../$(OUT)/min

$(OUT)/minc: | $(OUT)
	cd minc && go build -o ../$(OUT)/minc

$(OUT)/minvm: | $(OUT)
	$(MAKE) -C minvm

clean:
	-@rm -f $(OUT)/min $(OUT)/minc $(OUT)/minvm
	$(MAKE) -C minvm clean

# Individual targets

min: $(OUT)/min
minc: $(OUT)/minc
minvm: minvm-debug minvm-release

minvm-debug:
	$(MAKE) -C minvm debug

minvm-release:
	$(MAKE) -C minvm release

.PHONY: all clean min minc minvm-debug minvm-release
