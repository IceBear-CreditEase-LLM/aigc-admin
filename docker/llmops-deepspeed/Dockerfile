FROM nvidia/cuda:11.0.3-cudnn8-devel-ubuntu20.04
# FROM nvidia/cuda:12.2.2-runtime-ubuntu20.04

LABEL maintainer="timvan"

COPY ./sources.list /etc/apt/sources.list

ARG PYTHON_VERSION=3.10
ARG CONDA_PYTHON_VERSION=3
ARG CONDA_DIR=/opt/conda

# Instal basic utilities
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget unzip bzip2 sudo build-essential && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install miniconda
ENV PATH $CONDA_DIR/bin:$PATH
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    wget https://mirrors.tuna.tsinghua.edu.cn/anaconda/miniconda/Miniconda$CONDA_PYTHON_VERSION-latest-Linux-x86_64.sh -O /tmp/miniconda.sh && \
    echo 'export PATH=$CONDA_DIR/bin:$PATH' > /etc/profile.d/conda.sh && \
    /bin/bash /tmp/miniconda.sh -b -p $CONDA_DIR && \
    rm -rf /tmp/* && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN conda install -y python=$PYTHON_VERSION && \
    conda clean -t -i -p -y

RUN  pip install --upgrade pip && \
    rm -rf ~/.cache/pip

ENV CUDA_HOME=/usr/local/cuda
ENV CUDA_ROOT=$CUDA_HOME
ENV PATH=$PATH:$CUDA_ROOT/bin:$HOME/bin

COPY . /app

RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1
ENV CUDA_HOME=/usr/local/cuda

EXPOSE 6006
EXPOSE 5000

RUN chmod +x ./chat_main.sh

# 当Docker容器启动时执行您的脚本，使用定义的环境变量
CMD ["/app/chat_main.sh"]