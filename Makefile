install: 
	@cd apps/server \
    && mkdir -p app \
    && go build -o app/kvserver main.go \
    && cp resources/conf.toml app/conf.toml \
	&& cd ../../ && rm -rf dist &&  mv apps/server/app dist

