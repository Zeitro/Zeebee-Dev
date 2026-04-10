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

| Branch Type | Naming Convention                      | Purpose                                            |
|-------------|----------------------------------------|----------------------------------------------------|
| Master      | `master`                               | Stores all production code (frontend + backend)    |
| Staging     | `staging`                              | Stores all staging phase code (frontend + backend) |
| Feature     | `feat-{area}-{feature}`                | New feature development                            |
| Environment | `staging-{backend/frontend}-{feature}` | Staging preparation for specific area              |
| Environment | `prod-{backend/frontend}-{feature}`    | Direct deployment branch for specific area         |

**Key Concepts:**

- **`master`**: Contains all production-ready code for both frontend and backend
- **`staging`**: Contains all staging-phase code for both frontend and backend
- **`prod-{area}-{feature}`**: Production deployment branch for the specific area (backend/frontend), can be used to
  deploy directly

### Feature Development Process

1. **Create Feature Branch**
   - Fork from `master` branch
   - Name convention: `feat-{backend/frontend}-{feature-name}`
   - Examples:
      - `feat-backend-dashboard`
      - `feat-frontend-user-profile`

2. **Development**
   - Develop your feature in the feature branch
   - Follow coding standards in respective subdirectories
   - Test locally before pushing

3. **Merge to Staging**
   - Merge feature branch to `staging-{backend/frontend}-{feature}`
   - Then merge to `staging` branch (contains all staging code)
   - Test in staging environment

4. **Merge to Production**
   - After staging tests pass, merge to `prod-{backend/frontend}-{feature}` (for direct deployment)
   - Final merge to `master` branch (contains all production code)

### Merge Sequence

```
feat-{area}-{feature}
    ↓
staging-{area}-{feature}
    ↓
staging (contains all staging code)
    ↓ (after staging tests pass)
prod-{area}-{feature} (direct deployment branch)
    ↓
master (contains all production code)
```

### Deployment Notes

- Use `prod-{backend}-{feature}` or `prod-{frontend}-{feature}` branches for direct deployment
- `master` branch serves as the source of truth for all production code
- `staging` branch serves as the source of truth for all staging code
