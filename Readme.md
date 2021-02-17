# DFD

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Docker Image Version (latest by date)](https://img.shields.io/docker/v/devsdk/dfd-server)](https://hub.docker.com/repository/docker/devsdk/dfd-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/devsdk/DFD)](https://goreportcard.com/report/github.com/devsdk/DFD)

Front-End:


[https://devsdk.net/dfd](https://devsdk.net/dfd)

Repository : [https://github.com/DevSDK/DFD-WEB](https://github.com/DevSDK/DFD-WEB)

Stacks : React, Redux, React-Bootstrap, ApexChart, TypeScript, AXIOS

  
This web service provide League of Legends histories and statistics or crew status for our small group called 'fv".

---

DFD는 롤팟 'fv'을 위해 League of Legends 승률과 전적을 제공합니다.

![Untitled](https://user-images.githubusercontent.com/18409763/106856161-6d38cf80-6701-11eb-8f96-a285b380b327.png)

![Untitled 1](https://user-images.githubusercontent.com/18409763/106856193-76c23780-6701-11eb-8391-c1c26bf6ca95.png)

## Main Feature

- Game frequency (Like github)
- Win ratio histories chart
- Win and Total game count
- Game Histories List
- Crew Status (i.e. Today I'll rest)
- Image server

---

### 메인 기능

- 게임 빈도 (깃헙의 잔디와 같은 그것)
- 게임모드별 승률 차트
- 총합 승률 및 게임수
- 게임 전적
- 롤팟 게임 있는지 여부
- 이미지 서버

### Getting Started

Setup mongodb and redis server.

You could install with docker:

---

몽고디비와 Redis 서버가 필요합니다.

도커를 이용해 설치할 수 있습니다.

- [Mongo](https://hub.docker.com/_/mongo)
- [Redis](https://hub.docker.com/_/redis/)

You must need setup environment variables. 

---

환경변수 설정이 반드시 필요합니다.

| Name                 | example                              | description                                                                   |
|----------------------|--------------------------------------|-------------------------------------------------------------------------------|
| DB_LOCATION          | localhost:27017                      | Mongo DB location                                                             |
| DB_AUTH_ID           | root                                 | Mongo DB Auth Id                                                              |
| DB_AUTH_PASSWORD     | admin                                | Mongo DB Auth Password                                                        |
| REDIS_LOCATION       | localhost:6379                       | Redis Location                                                                |
| REDIS_PASSWORD       | 1234                                 | Redis auth password                                                           |
| SERVER_URI           | localhost                            | server location                                                               |
| BASE_URL             | /                                    | Base URL for cookie setting                                                   |
| REDIRECT_URL         | http://localhost:3000/auth           | Redirect URL after oauth2 success                                             |
| DISCORD_API_BASE     | https://discord.com/api/v6           | Discord API Location                                                          |
| DISCORD_CLIENT_ID    | 12345678910211                       | Discord API client id                                                         |
| DISCORD_REDIRECT_URI | http://localhost:18020/auth/redirect | Discord Oauth2 redirect location                                              |
| DISCORD_SCOPES       | identify email connections           | Discord access scopes                                                         |
| DISCORD_SECRET_ID    | asfipurofu9dias9c891 ....            | Discord secret code                                                           |
| DFD_SECRET_CODE      | dfoi1u20cvac801720d7cacs ....        | Secret code for DFD JWT authentication secret code (Recommended Random Hash)  |
| RIOT_API_URI         | https://kr.api.riotgames.com         | Riot API Location                                                             |


**1. Using docker** 
[https://hub.docker.com/repository/docker/devsdk/dfd-server](https://hub.docker.com/repository/docker/devsdk/dfd-server)

**2. Install all dependencies & Run**

```
git clone https://github.com/DevSDK/DFD.git
cd DFD
go mod download
go run main.go
```

### API Docs

[https://devsdk.net/api/dfd/docs/index.html](https://devsdk.net/api/dfd/docs/index.html)

### Stacks

![Untitled 2](https://user-images.githubusercontent.com/18409763/106857409-6317d080-6703-11eb-8004-efc43f436f91.png)



### Tasks
- [x]  Implement Endpoints
- [x]  Documentation
- [x]  Zero Downtime Deploy
- [x]  Front-End
- [ ]  Test
- [ ]  CI/CD

### For Admin

If you are admin, you could access with authentication: 

- [Mongo Express](https://devsdk.net/server/mongo/dashboard/)
- [Kubernetes Dashboard](https://devsdk.net/server/kube/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/pod?namespace=default)