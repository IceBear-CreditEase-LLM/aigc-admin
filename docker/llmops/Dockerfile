#FROM nvidia/cuda:11.0.3-cudnn8-devel-ubuntu20.04
FROM reg.creditease.corp/basic/other:nvidia-cuda-12.2.2-runtime-ubuntu20.04

LABEL maintainer="congwang@creditease.cn"

COPY ./sources.list /etc/apt/sources.list

ARG PYTHON_VERSION=3.10

ARG DEBIAN_FRONTEND=noninteractive

ARG http_proxy=http://10.170.32.94:7890
ARG https_proxy=http://10.170.32.94:7890
ARG all_proxy=http://10.170.32.94:7890
ARG no_proxy="10.170.*,10.181.*,10.182.*,.idc,.corp"

# Instal basic utilities
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget unzip bzip2 sudo build-essential jq vim wget git

RUN apt-get install -y apt-utils tzdata curl

RUN apt-get install -y cuda-toolkit --fix-missing

RUN apt-get install -y --no-install-recommends \
    software-properties-common

# 添加Deadsnakes PPA
RUN add-apt-repository ppa:deadsnakes/ppa
# 安装Python 3.10
RUN apt-get install -y --no-install-recommends \
    python3.10 \
    python3.10-venv \
    python3.10-distutils \
    python3.10-dev

RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
RUN python3 get-pip.py

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN pip install --upgrade pip && \
    rm -rf ~/.cache/pip

RUN pip3 install packaging torch==2.1.0
RUN pip3 install flash-attn==2.0.9 --no-build-isolation
RUN pip3 install fschat[model_worker,train] pydantic==1.10.13 optimum auto-gptq exllamav2 jsonlines
RUN pip3 install accelerate>=0.21 peft sentencepiece protobuf einops wandb anthropic>=0.3 ray
RUN pip3 install vllm==0.3.1

WORKDIR /app

COPY ./requirements.txt /app/requirements.txt

RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY . /app
