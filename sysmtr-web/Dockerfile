FROM node:20-alpine AS build
WORKDIR /app

COPY sysmtr-web/package*.json ./

RUN npm install --legacy-peer-deps

COPY sysmtr-web/ .

RUN npm run build -- --configuration production

FROM nginx:stable-alpine

COPY --from=build /app/dist/sysmtr-web/browser /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

