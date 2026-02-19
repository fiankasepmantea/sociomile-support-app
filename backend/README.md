1️⃣ README — FINAL SUBMISSION TEMPLATE
# Support Conversation System — Backend Service

## Overview
Backend REST API untuk sistem support conversation multi-tenant dengan fitur:

- Authentication (JWT)
- Multi-tenant isolation
- Conversation & messaging
- Customer management
- Ticketing
- Webhook ingestion
- Activity audit log

Tech stack:
- Golang (Gin)
- PostgreSQL (GORM)
- JWT auth
- REST API design

---

## Features

- Multi-tenant middleware (header based tenant resolution)
- JWT token authentication
- Conversation lifecycle (open → assigned → closed)
- Agent assignment with validation
- Message send & receive
- Webhook channel receiver
- Customer management
- Ticket creation
- Activity log trail

---

## Project Structure

cmd/
internal/
handler/
service/
repository/
middleware/
model/
dto/
pkg/database/



Layered architecture:

Handler → Service → Repository → Database

---

## Environment Variables

APP_ENV=development
SERVER_PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/support
JWT_SECRET=secret
JWT_EXPIRY_HOURS=24

---

## Run

go mod tidy
go run main.go

Health check:


Health check:


---

## Authentication

Header required:

Authorization: Bearer <token>
X-Tenant-ID: <tenant> atau <id>


---

## API Base


---

## Endpoint Coverage

Auth ✔  
Conversations ✔  
Messages ✔  
Customers ✔  
Tickets ✔  
Webhook ✔  

---

## Security

- Tenant isolation enforced at repository layer
- Role validation for agent actions
- JWT signed tokens
- Activity audit log

2️⃣ ARSITEKTUR DIAGRAM
                 ┌───────────────┐
                 │   Frontend    │
                 │ Dashboard UI  │
                 └───────┬───────┘
                         │ REST
                         ▼
                ┌──────────────────┐
                │   Gin Router     │
                └───────┬──────────┘
                        │
        ┌───────────────┼────────────────┐
        ▼               ▼                ▼
   Auth Middleware   Tenant MW       CORS MW
        │               │
        ▼               ▼
              Handler Layer
   ┌────────────┬────────────┬────────────┐
   │ Auth       │ Conversation│ Message   │
   └────┬───────┴─────┬───────┴─────┬─────┘
        ▼             ▼             ▼
               Service Layer
        Business rules & validation
        │
        ▼
           Repository Layer
      (tenant filtered queries)
        │
        ▼
             PostgreSQL


3️⃣ DESIGN DECISIONS
Layered Architecture
- Memisahkan HTTP, business logic, dan data access
- Mudah test & maintain

Multi-Tenant via Context
- Tenant injected via middleware
- Repository auto filter tenant_id

JWT Stateless Auth
- Tidak perlu session store
- Cocok untuk horizontal scaling

Service Validation
- Semua rule kritikal di service:
  - role check
  - assignment check
  - status transition

Audit Log
- Semua action penting tercatat
- Mendukung trace & compliance

Soft coupling Repo
- Repository interface pattern
- Mudah ganti DB

4️⃣ SWAGGER (OPENAPI 3.0)
openapi: 3.0.0
info:
  title: Support Conversation API
  version: 1.0.0

servers:
  - url: http://localhost:8080

paths:

  /api/auth/login:
    post:
      summary: Login
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email,password]
              properties:
                email: { type: string }
                password: { type: string }
      responses:
        200: { description: OK }

  /api/conversations:
    get:
      summary: List conversations
      security: [ bearerAuth: [] ]
      responses:
        200: { description: OK }

  /api/conversations/{id}:
    get:
      summary: Conversation detail
      security: [ bearerAuth: [] ]

  /api/conversations/{id}/assign:
    patch:
      summary: Assign agent

  /api/conversations/{id}/status:
    patch:
      summary: Update status

  /api/conversations/{id}/messages:
    get:
      summary: List messages
    post:
      summary: Send message

  /api/customers:
    get:
      summary: List customers
    post:
      summary: Create customer

  /api/customers/{id}:
    get:
      summary: Customer detail

  /api/tickets:
    post:
      summary: Create ticket

  /api/webhooks/channel:
    post:
      summary: Receive webhook

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer

5️⃣ POSTMAN COLLECTION (IMPORT JSON)
{
  "info": { "name": "Support API", "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" },
  "item": [

    {
      "name": "Login",
      "request": {
        "method": "POST",
        "url": "{{base}}/api/auth/login",
        "body": {
          "mode": "raw",
          "raw": "{ \"email\":\"admin@test.com\", \"password\":\"123456\" }"
        }
      }
    },

    {
      "name": "List Conversations",
      "request": {
        "method": "GET",
        "url": "{{base}}/api/conversations"
      }
    },

    {
      "name": "Send Message",
      "request": {
        "method": "POST",
        "url": "{{base}}/api/conversations/1/messages",
        "body": { "mode": "raw", "raw": "{ \"body\":\"hello\" }" }
      }
    }

  ]
}


6️⃣ CURL TEST SCRIPT — ALL ENDPOINTS
BASE=http://localhost:8080
TENANT=tenant1

echo LOGIN
TOKEN=$(curl -s -X POST $BASE/api/auth/login \
 -H "X-Tenant-ID: $TENANT" \
 -H "Content-Type: application/json" \
 -d '{"email":"admin@test.com","password":"123456"}' | jq -r .token)

AUTH="-H Authorization: Bearer $TOKEN -H X-Tenant-ID: $TENANT"

echo LIST CONVERSATIONS
curl $BASE/api/conversations $AUTH

echo GET DETAIL
curl $BASE/api/conversations/1 $AUTH

echo ASSIGN
curl -X PATCH $BASE/api/conversations/1/assign \
 $AUTH -d '{"agent_id":2}'

echo STATUS
curl -X PATCH $BASE/api/conversations/1/status \
 $AUTH -d '{"status":"closed"}'

echo LIST MSG
curl $BASE/api/conversations/1/messages $AUTH

echo SEND MSG
curl -X POST $BASE/api/conversations/1/messages \
 $AUTH -d '{"body":"reply"}'

echo CUSTOMERS
curl $BASE/api/customers $AUTH

echo CREATE CUSTOMER
curl -X POST $BASE/api/customers \
 $AUTH -d '{"external_id":"ext1","name":"John"}'

echo WEBHOOK
curl -X POST $BASE/api/webhooks/channel \
 -H "X-Tenant-ID: $TENANT" \
 -d '{"external_customer_id":"ext1","message_body":"hi"}'

echo TICKET
curl -X POST $BASE/api/tickets \
 $AUTH -d '{"title":"Issue","conversation_id":1}'


7️⃣ SUBMISSION CHECKLIST
✔ Multi-tenant
✔ JWT auth
✔ REST style
✔ Conversation lifecycle
✔ Message handling
✔ Webhook ingestion
✔ Customer CRUD
✔ Ticket create
✔ Activity log
✔ PostgreSQL
✔ Swagger spec
✔ Postman collection
✔ Curl tests
✔ README
✔ Architecture diagram
✔ Design decisions


## ⚡ Redis Optimization

Redis is used to improve performance and reliability:

- Webhook rate limiting (anti spam)
- JWT blacklist for logout
- Conversation list caching
- Conversation detail caching

Design principle:
- Redis optional
- Graceful fallback if unavailable
- Short TTL cache (60s) to avoid stale data

This ensures horizontal scalability and performance under load.

