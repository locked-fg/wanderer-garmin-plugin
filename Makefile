# Garmin Plugin Makefile
# Builds the WASM plugin for wanderer

PLUGIN_ID := garmin
PLUGIN_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DIST_DIR := $(PLUGIN_DIR)/dist/$(PLUGIN_ID)

# TinyGo settings (same version as wanderer CI)
TINYGO_VERSION ?= 0.41.1
TINYGO ?= tinygo

# Plugin files
PLUGIN_FILES := plugin.json plugin.wasm icon.svg

.PHONY: all build manifest install-local clean package

all: build

# Validate and normalize the plugin manifest into dist/
manifest:
	@mkdir -p $(DIST_DIR)
	cd $(PLUGIN_DIR) && go run github.com/open-wanderer/wanderer/plugins/sdk/cmd/manifestcheck plugin.json > $(DIST_DIR)/plugin.json
	cp $(PLUGIN_DIR)/assets/icon.svg $(DIST_DIR)/icon.svg

# Build the WASM plugin
build: manifest
	@echo "Building Garmin plugin..."
	cd $(PLUGIN_DIR) && $(TINYGO) build -target=wasi -scheduler=none -no-debug -o $(DIST_DIR)/plugin.wasm .
	@echo "Built: $(DIST_DIR)/plugin.wasm"

# Install to local wanderer data directory (repo layout: plugins/garmin -> data/plugins/garmin)
install-local: build
	@echo "Installing to local wanderer..."
	@mkdir -p $(PLUGIN_DIR)/../../data/plugins/$(PLUGIN_ID)
	cp $(DIST_DIR)/plugin.json $(PLUGIN_DIR)/../../data/plugins/$(PLUGIN_ID)/
	cp $(DIST_DIR)/plugin.wasm $(PLUGIN_DIR)/../../data/plugins/$(PLUGIN_ID)/
	cp $(DIST_DIR)/icon.svg $(PLUGIN_DIR)/../../data/plugins/$(PLUGIN_ID)/
	@echo "Plugin installed to: $(PLUGIN_DIR)/../../data/plugins/$(PLUGIN_ID)"

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR)

# Package the plugin for distribution
package: build
	@echo "Packaging Garmin plugin..."
	cd $(DIST_DIR) && tar czf ../wanderer-plugin-$(PLUGIN_ID).tar.gz $(PLUGIN_FILES)
	@echo "Package created: $(DIST_DIR)/../wanderer-plugin-$(PLUGIN_ID).tar.gz"
