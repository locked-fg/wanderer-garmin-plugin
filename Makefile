# Garmin Plugin Makefile
# Builds the WASM plugin for wanderer

PLUGIN_ID := garmin
PLUGIN_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DIST_DIR := dist/$(PLUGIN_ID)
WANDERER_DATA_DIR ?= ../data
LOCAL_PLUGIN_DIR := $(WANDERER_DATA_DIR)/plugins/$(PLUGIN_ID)

# TinyGo settings
TINYGO_VERSION ?= 0.44.0
TINYGO ?= tinygo
TINYGO_TARGET ?= wasi
TINYGO_FLAGS ?= -scheduler=none -no-debug

# Plugin files
PLUGIN_FILES := plugin.json plugin.wasm icon.svg

.PHONY: all build install clean

all: build

# Build the WASM plugin
build: plugin.wasm

GO_FILES := $(wildcard *.go) $(wildcard internal/*/*.go)

plugin.wasm: $(GO_FILES) go.mod go.sum
	@echo "Building Garmin plugin..."
	cd $(PLUGIN_DIR) && mkdir -p $(DIST_DIR) && $(TINYGO) build $(TINYGO_FLAGS) -target $(TINYGO_TARGET) -o $(DIST_DIR)/plugin.wasm .
	@echo "Built: $(DIST_DIR)/plugin.wasm"

# Install the plugin locally
install: build
	@echo "Installing Garmin plugin..."
	cd $(PLUGIN_DIR) && mkdir -p $(DIST_DIR)
	cd $(PLUGIN_DIR) && cp plugin.json $(DIST_DIR)/
	cd $(PLUGIN_DIR) && cp assets/icon.svg $(DIST_DIR)/icon.svg
	@echo "Plugin installed to: $(DIST_DIR)"

# Clean build artifacts
clean:
	cd $(PLUGIN_DIR) && rm -rf $(DIST_DIR) && rm -f plugin.wasm

# Package the plugin for distribution
package: build
	@echo "Packaging Garmin plugin..."
	cd $(PLUGIN_DIR) && mkdir -p $(DIST_DIR) && cp plugin.json $(DIST_DIR)/ && cp assets/icon.svg $(DIST_DIR)/icon.svg
	cd $(PLUGIN_DIR)/$(DIST_DIR) && tar czf ../../wanderer-plugin-$(PLUGIN_ID).tar.gz $(PLUGIN_FILES)
	@echo "Package created: ../../wanderer-plugin-$(PLUGIN_ID).tar.gz"

# Install to local wanderer data directory
install-local: install
	@echo "Installing to local wanderer..."
	cd $(PLUGIN_DIR) && mkdir -p "$(LOCAL_PLUGIN_DIR)" && cp $(DIST_DIR)/plugin.json $(DIST_DIR)/plugin.wasm $(DIST_DIR)/icon.svg "$(LOCAL_PLUGIN_DIR)/"
	@echo "Plugin installed to: $(LOCAL_PLUGIN_DIR)"

# Build all plugins
plugins-build:
	cd .. && make plugins-build

# Install all plugins locally
plugins-install-local:
	cd .. && make plugins-install-local
