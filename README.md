# Zeebee-Dev

Private lending mortgage platform repository

## Overview

Zeebee is a comprehensive private lending mortgage platform built with a microservices architecture. This monorepo
contains both backend and frontend components.

## Project Structure

```
Zeebee-Dev/
├── backend/      # Go + go-zero microservices
├── frontend/     # Frontend application
```

## Getting Started

### Prerequisites

- Go (for backend development)
- Node.js (for frontend development)
- Git

### Initial Setup

1. Clone the repository:
   ```bash
   git clone git@github.com:Zeitro/Zeebee-Dev.git
   cd Zeebee-Dev
   ```

2. For backend development, see [backend/README.md](backend/README.md)
3. For frontend development, see [frontend/README.md](frontend/README.md)

## Development Workflow

### Branch Strategy

We follow a Gitflow-inspired branching model:

| Branch Type | Naming Convention                                    | Purpose                                            |
|-------------|------------------------------------------------------|----------------------------------------------------|
| Master      | `master`                                             | Stores all production code (frontend + backend)    |
| Staging     | `staging`                                            | Stores all staging phase code (frontend + backend) |
| Feature     | `feat-{borrower/lender/admin}-{feature}`             | New feature development                            |
| Environment | `staging-{borrower/lender/admin}-{backend/frontend}` | Staging preparation for specific area              |
| Environment | `prod-{borrower/lender/admin}-{backend/frontend}`    | Direct deployment branch for specific area         |

**Key Concepts:**

- **`master`**: Contains all production-ready code for both frontend and backend
- **`staging`**: Contains all staging-phase code for both frontend and backend
- **Business Areas**: `borrower`, `lender` , `admin`
- **Tech Areas**: `backend` , `frontend`

### Feature Development Process

1. **Create Feature Branch**
   - Fork from `master` branch
   - Name convention: `feat-{business-area}-{feature-name}`
   - Examples:
       - `feat-borrower-user-profile`
       - `feat-lender-loan-approval`
       - `feat-admin-role-management`

2. **Development**
   - Develop your feature in the feature branch
   - Follow coding standards in respective subdirectories
   - Test locally before pushing

3. **Merge to Staging**
    - Merge feature branch to `staging-{business-area}-{backend/frontend}`
   - Then merge to `staging` branch (contains all staging code)
   - Test in staging environment

4. **Merge to Production**
    - After staging tests pass, merge to `prod-{borrower/lender/admin}-{backend/frontend}` (for direct deployment)
   - Final merge to `master` branch (contains all production code)

### Merge Sequence

```
feat-{business-area}-{feature}
    ↓
staging-{business-area}-{backend/frontend}

-------------------------------------

feat-{business-area}-{feature}
    ↓
prod-{borrower/lender/admin}-{backend/frontend} (direct deployment branch)

--------------------------------------

staging-{business-area}-{backend/frontend}
    ↓
staging (contains all staging code)
  
--------------------------------------    
    
prod-{borrower/lender/admin}-{backend/frontend} (direct deployment branch)
    ↓
master (contains all production code)

```

### Deployment Branch Examples

| Purpose                   | Branch Name                 |
|---------------------------|-----------------------------|
| Borrower Backend Staging  | `staging-borrower-backend`  |
| Borrower Frontend Staging | `staging-borrower-frontend` |
| Lender Backend Staging    | `staging-lender-backend`    |
| Lender Frontend Staging   | `staging-lender-frontend`   |
| Admin Backend Staging     | `staging-admin-backend`     |
| Admin Frontend Staging    | `staging-admin-frontend`    |
| Borrower Backend Prod     | `prod-borrower-backend`     |
| Borrower Frontend Prod    | `prod-borrower-frontend`    |
| Lender Backend Prod       | `prod-lender-backend`       |
| Lender Frontend Prod      | `prod-lender-frontend`      |
| Admin Backend Prod        | `prod-admin--backend`       |
| Admin Frontend Prod       | `prod-admin--frontend`      |

### Deployment Notes

- Use `prod-{business-area}--{backend/frontend}` branches for direct deployment
- `master` branch serves as the source of truth for all production code
- `staging` branch serves as the source of truth for all staging code
