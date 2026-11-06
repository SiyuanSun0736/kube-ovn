CLANG ?= clang
LLC ?= llc
ARCH := $(shell uname -m | sed 's/x86_64/x86/' | sed 's/aarch64/arm64/')

BPF_CFLAGS ?= -O2 -g -Wall -target bpf \
	-D__TARGET_ARCH_$(ARCH) \
	-I/usr/include \
	-I/usr/include/$(shell uname -m)-linux-gnu

OUTPUT_DIR ?= .
TARGET := $(OUTPUT_DIR)/tc_eip.o

all: $(TARGET)

$(TARGET): tc_eip.c
	$(CLANG) $(BPF_CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET)

.PHONY: all clean
