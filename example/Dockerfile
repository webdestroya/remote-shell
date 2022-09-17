FROM httpd:2.4

COPY index.html /usr/local/apache2/htdocs/index.html

COPY --from=ghcr.io/webdestroya/remote-shell:latest /cloud87 /cloud87

