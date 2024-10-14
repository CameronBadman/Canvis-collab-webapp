# Build stage
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
ARG VITE_API_URL
ARG VITE_AUTH_URL
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_AUTH_URL=${VITE_AUTH_URL}
RUN echo "VITE_API_URL: $VITE_API_URL"
RUN echo "VITE_AUTH_URL: $VITE_AUTH_URL"
RUN npm run build
RUN ls -la /app/dist

# Production stage
FROM nginx:stable-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3002
CMD ["nginx", "-g", "daemon off;"]