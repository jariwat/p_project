# Bundle OpenAPI specifications
bundle:
	@if [ -z "$(SPECS)" ]; then \
		echo "No OpenAPI specifications found"; \
		exit 1; \
	fi
	@for spec in $(SPECS); do \
		service_dir=$$(dirname $$spec); \
		echo "Bundling $$service_dir..."; \
		redocly bundle $$spec --output $$service_dir/openapi_bundle.yml; \
		redocly bundle $$spec --output $$service_dir/openapi_bundle.json; \
	done