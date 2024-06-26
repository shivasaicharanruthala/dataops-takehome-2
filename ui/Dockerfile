# Pull the base docker image of python with tag 3.11
FROM python:3.11

# Upgrade pip package
RUN pip install --upgrade pip

# Change the working dir inside the container - cd /ui
WORKDIR /ui

# Copy main.py as source cod and req.txt as dependency
COPY  main.py requirements.txt ./

# Install the dependency
RUN pip install -r requirements.txt

### Container Env same as local at this point
EXPOSE 8090

CMD ["streamlit", "run", "main.py", "--server.port", "8090"]
