# Use an official Go runtime as a parent image
FROM golang:1.16-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Set environment variables
ENV DB_URI=postgres://default:CYJ7uvnw3Xma@ep-jolly-glade-04413070.us-east-1.postgres.vercel-storage.com:5432/verceldb?sslmode=require
ENV DB_USER=if0_34940045
ENV DB_PASSWORD=password
ENV DB_HOST=localhost
ENV DB_PORT=5433
ENV DB_NAME=alutamarket-db
ENV SECRET_KEY=fZYuSCOsDCl5t0150eb40MRo-1R3z7tpt0NrvlRpiX0IZYuSCOsDCl5t0150eb40MRo-1R3z7tpt0NrvlRpiX0I
ENV DOMAIN=localhost
ENV SENDGRID_KEY=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
ENV SENDER_EMAIL=opeyemifolajimi13@gmail.com
ENV TWILIO_ACCOUNT_SID=ACba0dac0e4926b219859aef814770fedf
ENV TWILIO_AUTH_TOKEN=148f65a5f294d7f7c663e182bc484978
ENV PORT=8082
ENV RAPID_API_KEY=97e0debb1dmshe760970bba1a576p1353a3jsn11f346494b39
ENV ACCESS_TOKEN=
ENV FLW_SECRET_KEY=FLWSECK_TEST-2697e7a01d28ec88bca63059be559903-X
ENV FLW_PUBLIC_KEY=FLWPUBK_TEST-fe4c090474e4141515a91edf955dac6a-X
ENV FLW_SECRET_HASH=FolAlutagengz
ENV PAYSTACK_SECRET_KEY=
ENV PAYSTACK_PUBLIC_KEY=
ENV PAYSTACK_SECRET_HASH=
ENV AWS_SECRET_KEY=gT39zemlHTdHj9vC4DKmQ0STMdw7q0MN/5B6kCks
ENV AWS_ACCESS_KEY=AKIA5VKPDOR5GWYLFZIJ

# Build the Go application
RUN go build -o server .

# Run the server when the container starts
CMD ["./server"]
