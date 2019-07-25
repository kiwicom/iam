FROM nginx:1.17.2-alpine

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./mocks ./mocks

VOLUME ./mocks
