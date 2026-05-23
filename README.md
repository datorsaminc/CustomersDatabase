# CustomersDatabase

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

A web-based customer database management system built with Go and Bootstrap 5. Originally migrated from Excel XML data to a modern JSON-backed REST API, providing full CRUD operations and real-time search capabilities through a responsive frontend.

## Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Data Model](#data-model)
- [Development](#development)

## Features

- **Full-text Search** — Real-time search across all customer fields (company, name, email, phone, address, etc.)
- **CRUD Operations** — Create, Read, Update, and Delete customer records via a modal-based UI
- **Pagination** — Browse large datasets efficiently with 50 records per page
- **Sorting** — Click column headers to sort by company, contact name, email, phone, or city
- **License Tracking** — Per-customer license counts with aggregate totals displayed on the dashboard
- **Form Validation** — Client-side validation for required fields and email format
- **Responsive Design** — Fully functional on desktop and mobile devices

## Technology Stack

### Backend (Go)

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| Router | [chi/v5](https://github.com/go-chi/chi) |
| Storage | JSON file-based (`data/customers.json`, `data/licenses.json`) |
| IDs | [google/uuid](https://github.com/google/uuid) |

### Frontend

| Component | Technology |
|-----------|------------|
| Framework | Bootstrap 5.3.2 |
| JavaScript | jQuery 3.7.1 |
| Architecture | Modular JS (API layer, table rendering, search, form handling) |

## Project Structure

```
ClientsWeb/
├── Kunder.xml                          # Original Excel data (source)
├── README.md                           # This file
├── plans/
│   └── architecture_plan.md            # Architecture documentation
├── backend/
│   ├── go.mod, go.sum                  # Go module definition
│   ├── config.go                       # Configuration loader
│   ├── config.json                     # Server configuration (port)
│   ├── main.go                         # Application entry point
│   ├── internal/
│   │   ├── handler/
│   │   │   └── customer_handler.go     # HTTP request handlers
│   │   ├── model/
│   │   │   ├── customer.go             # Customer data model
│   │   │   └── license.go              # License data model
│   │   ├── search/
│   │   │   └── searcher.go             # Full-text search logic
│   │   ├── service/
│   │   │   └── customer_service.go     # Business logic layer
│   │   └── store/
│   │       ├── json_store.go           # Generic JSON file storage
│   │       └── license_store.go        # License-specific storage
│   └── data/
│       ├── customers.json              # Customer database (generated)
│       └── licenses.json               # License database (generated)
└── frontend/
    ├── index.html                      # Main application page
    ├── userhelp.html                   # User help documentation
    ├── css/
    │   └── style.css                   # Custom styles
    └── js/
        ├── app.js                      # Application initialization
        ├── api.js                      # API communication layer
        ├── search.js                   # Search functionality
        ├── customerForm.js             # Form handling & validation
        └── customerTable.js            # Table rendering, sorting & pagination
```

## Getting Started

### Prerequisites

- **Go 1.25+** installed on your system ([download](https://go.dev/dl/))

### Quick Start

1. **Build and start the server:**

   ```bash
   cd backend
   go build -o server.exe .
   ./server.exe
   ```

2. **Open your browser** and navigate to [http://localhost:8080](http://localhost:8080)

The application serves both the REST API (at `/api/customers`) and the frontend static files from a single process.

### Running Without Building

You can also run directly with `go run`:

```bash
cd backend
go run .
```

## Configuration

Edit [`backend/config.json`](backend/config.json) to change server settings:

```json
{
  "port": "8080"
}
```

| Setting | Description | Default |
|---------|-------------|---------|
| `port`  | HTTP listen port | `8080` |

## API Endpoints

### Customers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/customers?search=query&page=1&limit=50` | List customers with optional search and pagination |
| `GET` | `/api/customers/:id` | Get a single customer by ID |
| `POST` | `/api/customers` | Create a new customer |
| `PUT` | `/api/customers/:id` | Update an existing customer |
| `DELETE` | `/api/customers/:id` | Delete a customer |

### Licenses

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/licenses/stats` | Get aggregate license statistics |

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `search`  | string | — | Full-text search query (matches across all text fields) |
| `page`    | int | `1` | Page number for pagination |
| `limit`   | int | `50` | Maximum results per page |

### Example Requests

```bash
# Search for customers matching "ABB"
curl "http://localhost:8080/api/customers?search=ABB&limit=5"

# Get a specific customer
curl "http://localhost:8080/api/customers/<customer-id>"

# Create a new customer
curl -X POST "http://localhost:8080/api/customers" \
  -H "Content-Type: application/json" \
  -d '{"company":"Acme Corp","name1":"John Doe","email":"john@acme.com"}'
```

## Data Model

### Customer

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string (UUID) | auto | Unique identifier |
| `programVersion` | string | no | Software version(s) owned |
| `deliveryDate` | string | no | Delivery date information |
| `name1` | string | **yes** | Primary contact name |
| `name2` | string | no | Secondary contact name |
| `company` | string | **yes** | Company name |
| `visitAddress` | string | no | Physical visit address |
| `mailingAddress` | string | no | Mailing address |
| `postalCodeCity` | string | no | Postal code and city |
| `landlinePhone` | string | no | Landline phone number |
| `mobilePhone` | string | no | Mobile phone number |
| `faxNumber` | string | no | Fax number |
| `email` | string | no | Email address |
| `comments` | string | no | Additional notes |
| `createdAt` | string (RFC3339) | auto | Record creation timestamp |
| `updatedAt` | string (RFC3339) | auto | Last update timestamp |

### License

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Customer ID (foreign key reference) |
| `licences` | int | Number of licenses for this customer |

## Development

### Architecture

The backend follows a clean layered architecture:

```
HTTP Handler → Service Layer → Store Layer → JSON File
```

- **Handler** — Parses HTTP requests, validates input, returns responses
- **Service** — Business logic (CRUD operations, search coordination)
- **Store** — Data persistence abstraction (currently JSON file-based)

This separation allows swapping the storage backend (e.g., to a database) without changing handler or service code.

### Adding New Features

1. Define the data model in `internal/model/`
2. Implement storage in `internal/store/`
3. Add business logic in `internal/service/`
4. Expose HTTP endpoints in `internal/handler/`
5. Register routes in `main.go`

---

**Built with Go & Bootstrap 5**
