# Garmin Plugin Makefile
# Builds the WASM plugin for wanderer

PLUGIN_ID := garmin
PLUGIN_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DIST_DIR := $(PLUGIN_DIR)/dist/$(PLUGIN_ID)

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

plugin.wasm: main.go garmin.go types.go
	@echo "Building Garmin plugin..."
	@mkdir -p $(DIST_DIR)
	cd $(PLUGIN_DIR) && $(TINYGO) build $(TINYGO_FLAGS) -target $(TINYGO_TARGET) -o $(DIST_DIR)/plugin.wasm .
	@echo "Built: $(DIST_DIR)/plugin.wasm"

# Install the plugin locally
install: build
	@echo "Installing Garmin plugin..."
	@mkdir -p $(DIST_DIR)
	cp $(PLUGIN_DIR)/plugin.json $(DIST_DIR)/
	cp $(PLUGIN_DIR)/assets/icon.svg $(DIST_DIR)/icon.svg
	@echo "Plugin installed to: $(DIST_DIR)"

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR)
	cd $(PLUGIN_DIR) && rm -f plugin.wasm

# Package the plugin for distribution
package: build
	@echo "Packaging Garmin plugin..."
	@mkdir -p $(DIST_DIR)
	cp $(PLUGIN_DIR)/plugin.json $(DIST_DIR)/
	cp $(PLUGIN_DIR)/assets/icon.svg $(DIST_DIR)/icon.svg
	cd $(DIST_DIR) && tar czf ../../wanderer-plugin-$(PLUGIN_ID).tar.gz $(PLUGIN_FILES)
	@echo "Package created: ../../wanderer-plugin-$(PLUGIN_ID).tar.gz"

# Install to local wanderer data directory
install-local: build
	@echo "Installing to local wanderer..."
	@mkdir -p ../../../data/plugins/$(PLUGIN_ID)
	cp $(DIST_DIR)/plugin.json ../../../data/plugins/$(PLUGIN_ID)/
	cp $(DIST_DIR)/plugin.wasm ../../../data/plugins/$(PLUGIN_ID)/
	cp $(PLUGIN_DIR)/assets/icon.svg ../../../data/plugins/$(PLUGIN_ID)/icon.svg
	@echo "Plugin installed to: ../../../data/plugins/$(PLUGIN_ID)"

# Build all plugins
plugins-build:
	cd .. && make plugins-build

# Install all plugins locally
plugins-install-local:
	cd .. && make plugins-install-local
