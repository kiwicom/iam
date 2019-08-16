FROM nginx:1.17.3-alpine

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./mocks ./mocks

VOLUME ./mocks
