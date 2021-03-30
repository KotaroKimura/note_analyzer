FROM golang:1.16

# Install packages
RUN \
  apt-get update && \
  apt-get upgrade -y && \
  apt-get install -y wget curl unzip apt-utils sudo

# Install python
RUN \
  apt-get install -y python3 python3-distutils && \
  curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py && \
  python3 get-pip.py && \
  pip install --upgrade pip

# Install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add && \
    echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable

# Install Selenium
RUN pip install selenium

# Install Chrome Driver
ADD https://chromedriver.storage.googleapis.com/89.0.4389.23/chromedriver_linux64.zip /opt/chrome/
RUN cd /opt/chrome/ && \
    unzip chromedriver_linux64.zip

# Install Font
RUN apt-get install -y fonts-migmix

# Set Japanese
RUN apt-get install -y locales
RUN locale-gen ja_JP.UTF-8

# Set Environment Variables.
ENV WORKSPACE /go/src
ENV LANG ja_JP.UTF-8
ENV PYTHONIOENCODIND utf_8
ENV PATH $PATH:/opt/chrome

# Install go packages
RUN go get golang.org/x/tools/cmd/goimports

# Project dir
RUN mkdir -p $WORKSPACE
WORKDIR $WORKSPACE
