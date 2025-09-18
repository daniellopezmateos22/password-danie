# --- build ---
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci || npm i
COPY . .
# URL de API para el build de Vite
ARG VITE_API_URL=http://api:8080
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# --- serve ---
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
