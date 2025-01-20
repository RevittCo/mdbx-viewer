FROM node:latest AS builder

WORKDIR /web
COPY . .

RUN npm install
RUN npm run build

# --------------------------------------------------------------------------------------------

FROM nginx:alpine

# Copy the html files from the builder stage to the nginx html directory
COPY --from=builder /web/dist /usr/share/nginx/html

# Copy the nginx config file to the nginx conf directory
COPY --from=builder /web/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]