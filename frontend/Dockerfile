FROM node:18 AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

# Copy the rest of the source code
COPY . .

RUN npm run build-only

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 5001

CMD ["nginx", "-g", "daemon off;"]
