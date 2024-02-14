# Dockerfile for MultLdr
FROM python

# Copy the challenge files
COPY /webapp /webapp

# Set the working directory
WORKDIR /webapp

# Install the challenge requirements
RUN pip3 install -r requirements.txt

# Expose the port 8080
EXPOSE 8080

# Run the command to start uWSGI
CMD ["uwsgi", "app.ini"]