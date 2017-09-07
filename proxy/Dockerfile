FROM nginx:alpine

RUN rm /etc/nginx/conf.d/*

COPY key.pem /etc/nginx/ssl/
COPY cert.pem /etc/nginx/ssl/

COPY proxy.conf /etc/nginx/conf.d/
