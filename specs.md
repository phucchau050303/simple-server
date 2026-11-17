# Specification: Go Registration System

## Overview
- **Purpose:** Create a Sign Up and Sign In system to learn and practice Golang, Temporal Workflow and gRPC.

## Goals
- User signs up with **email** and **password**
- System creates a pending registration and emails a 6-digit code
- User enters the code:
    - correct → account active, success page
    - wrong too many times → registration locked
    - lost / expired → user requests a new code
- Active user logs in, receives an access token, views profile
- Never verified → registration expires after 24h, cleaned up, email can register again

## Architecture
```mermaid
flowchart TD
    FE[Frontend] -->|HTTP| RS[Registration Server]
    
    subgraph Backend Services
        RS -->|gRPC| CS[Credential Service]
        RS -->|Task Queue / gRPC| TW[Temporal Worker]
    end
    
    CS -->|SQL| DB[(Postgres)]

    style FE fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    style RS fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style CS fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style TW fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style DB fill:#f3e5f5,stroke:#6a1b9a,stroke-width:2px
```

**Component List**
- **Registration Server:** Public HTTP API, UI, prospect lifecycle
- **Credential Service:** Internal gRPC, password hashing + verification
- **Temporal Worker:** Registration workflow, email, expiry

## Use Cases
| |	Use case | Rules |
| --- | --- | --- |
|UC1 | Register	Email + password policy validation.| No account enumeration. |
UC2	| Verify email	| 6-digit OTP, TTL 15m, max 5 attempts | 
UC3	| Resend OTP	| Invalidate previous OTP, rate-limited |
UC4	| Log in	| Active accounts only, JWT/PASETO |
UC5	| Get profile	| Auth middleware |
UC6	| Expire registrations	| Pending > 24h → expired, cleanup|

Prospect status: pending → verified → active / expired / locked

## Technical Design
### List of endpoints:
All endpoints accept and return JSON unless stated otherwise. The profile endpoint
requires a valid access token in the `Authorization` header.

| Method | Path | Authentication | Request body | Success response | Description |
| --- | --- | --- | --- | --- | --- |
| `POST` | `/signup` | None | `{ "email": "string", "password": "string" }` | `201 Created` | Create a pending registration and send an OTP. |
| `POST` | `/signin` | None | `{ "email": "string", "password": "string" }` | `200 OK` | Authenticate an active user and return an access token. |
| `POST` | `/otp/resend` | None | `{ "email": "string" }` | `202 Accepted` | Invalidate the previous OTP and send a new one. |
| `POST` | `/otp/verify` | None | `{ "email": "string", "otp": "string" }` | `200 OK` | Verify the OTP and activate the account. |
| `GET` | `/userprofile` | Bearer token | None | `200 OK` | Return the authenticated user's profile. |

### Database Schema:
Will be added later

## POC Flows

#### Sign up POC flow

```mermaid
sequenceDiagram
    actor User
    participant FE as Frontend
    participant RS as Registration Server
    participant CS as Credential Service
    participant DB as Postgres
    participant TW as Temporal Worker

    User->>FE: Enter email and password
    FE->>RS: POST /signup
    RS->>CS: Call credential service
    CS->>CS: Hash password
    CS->>DB: Save pending registration and hash
    RS->>TW: Start email workflow
    TW-->>User: Send 6-digit OTP
    RS-->>FE: 201 Created
    FE-->>User: Display OTP entry screen
```

#### OTP verification POC flow

```mermaid
flowchart TD
    A[User enters OTP] --> B[POST /otp/verify]
    B --> C{Registration status}
    C -->|Locked or expired| D[Return 400 Bad Request]
    C -->|Active OTP within 15 minutes| E[Transition pending → verified → active]
    E --> F[Return 200 OK]
    C -->|Incorrect or expired OTP| G[Increment attempts]
    G --> H{Attempts >= 5?}
    H -->|Yes| I[Set status to locked]
    I --> J[Return 400 Bad Request: Account locked]
    H -->|No| K[Return 400 Bad Request: Invalid code and attempts remaining]
```

#### Resend OTP POC flow

```mermaid
sequenceDiagram
    actor User
    participant FE as Frontend
    participant RS as Registration Server
    participant TW as Temporal Worker

    User->>FE: Request a new OTP
    FE->>RS: POST /otp/resend
    RS->>RS: Invalidate previous OTP
    RS->>TW: Generate and send a new OTP
    RS-->>FE: 202 Accepted
```

#### Sign in POC flow

```mermaid
sequenceDiagram
    actor User
    participant FE as Frontend
    participant RS as Registration Server
    participant CS as Credential Service
    participant DB as Postgres

    User->>FE: Enter email and password
    FE->>RS: POST /signin
    RS->>CS: Call service via gRPC
    CS->>DB: Verify password and hash (gRPC)
    CS-->>RS: Password result
    RS->>RS: Create signed JWT
    RS-->>FE: 200 OK with JWT
    FE->>FE: Store JWT
    FE-->>User: Redirect to profile
```

#### Get profile POC flow

```mermaid
sequenceDiagram
    actor User
    participant FE as Frontend
    participant RS as Registration Server
    participant CS as Credential Service
    participant DB as Postgres


    User->>FE: Open profile page
    FE->>RS: GET /userprofile with Authorization header
    RS->>RS: Validate JWT signature and expiry
    RS->>RS: Extract user ID/email from claims
    RS->>CS: Call credential service
    CS->>DB: Query user profile
    DB-->>RS: Return profile
    RS-->>FE: 200 OK with user profile
```
