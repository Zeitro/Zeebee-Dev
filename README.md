# Zeebee-Dev

Private lending mortgage platform monorepo.

## Structure

```
Zeebee-Dev/
├── backend/      # Go + go-zero microservices (borrower, lender, admin)
├── frontend/     # Frontend application
└── env/          # Shared environment configs
```

## Getting Started

```bash
git clone git@github.com:Zeitro/Zeebee-Dev.git
cd Zeebee-Dev

# Backend: see backend/README.md
# Frontend: see frontend/README.md
```

## Branch Strategy

Gitflow-inspired model with three business areas (`borrower`, `lender`, `admin`) and two tech areas (`backend`, `frontend`).

### Branch Types

| Type | Pattern | Purpose |
|------|---------|---------|
| Master | `master` | All production code |
| Staging | `staging` | All staging code |
| Feature | `feat-{area}-{name}` | New feature development |
| Staging Env | `staging-{area}-{tech}` | Staging preparation |
| Prod Env | `prod-{area}-{tech}` | Direct deployment |

**Examples**: `feat-borrower-user-auth`, `staging-lender-backend`, `prod-admin-frontend`

### Merge Flow

```
feat-{area}-{name}
    ↓
staging-{area}-{tech} → staging (all staging code)
    ↓ (after staging tests pass)
prod-{area}-{tech} → master (all production code)
```

### Deployment Branches

| Purpose | Branch |
|---------|--------|
| Borrower Backend Staging | `staging-borrower-backend` |
| Borrower Frontend Staging | `staging-borrower-frontend` |
| Lender Backend Staging | `staging-lender-backend` |
| Lender Frontend Staging | `staging-lender-frontend` |
| Admin Backend Staging | `staging-admin-backend` |
| Admin Frontend Staging | `staging-admin-frontend` |
| Borrower Backend Prod | `prod-borrower-backend` |
| Borrower Frontend Prod | `prod-borrower-frontend` |
| Lender Backend Prod | `prod-lender-backend` |
| Lender Frontend Prod | `prod-lender-frontend` |
| Admin Backend Prod | `prod-admin-backend` |
| Admin Frontend Prod | `prod-admin-frontend` |

---

*Project homepage: [github.com/Zeitro/Zeebee-Dev](https://github.com/Zeitro/Zeebee-Dev)*
