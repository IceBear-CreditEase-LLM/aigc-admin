package langchain

type Middleware func(Service) Service
