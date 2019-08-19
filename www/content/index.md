---
title: "Kiwi IAM"
date: 2019-08-20T11:56:33+02:00
draft: false
---

# Kiwi IAM

Kiwi IAM is a RBAC (Role based access control) app made for the needs of Kiwi. The aim is to simplify securing apps by providing all the relevant libraries and infrastructure.

**Features:**

- Serving a permission list of the user
- Caching Okta provider (Okta has a rate limit of 600 req/s)
- Serving Okta user profile for apps that need to use user data

# API service

API service is responsible for handling communication with Okta.

![Alt text](https://g.gravizo.com/svg?
@startuml;
== Happy path ==;
Router -> Controller : Request from client for test@test.com;
Controller -> Cache : Request test@test.com from Cache;
Cache -> Controller : Cache hit;
Controller -> Router : Return profile of test@test.com to client;
== Unhappy path ==;
Router -> Controller : Request from client for test@test.com;
Controller -> Cache : Request test@test.com from Cache;
Cache -> Controller : Profile not found;
Controller -> OktaAPI : Request test@test.com from Okta;
OktaAPI -> Controller : Return test@test.com to Controller;
Controller --> Cache : Send data from test@test.com to be saved in Cache;
Controller -> Router : Return data for test@test.com to client;
@enduml;
)
