# Testing Guide: Lender Organization Management

**Feature**: 003-lender-organization - Lender端组织管理功能
**Date**: 2025-01-09
**Purpose**: Manual testing guide for organization management APIs

## Prerequisites

1. **Database Running**: Ensure MySQL/PostgreSQL database is running and accessible
2. **Environment Configured**: Check `app/lender/etc/lender.yaml` for database connection
3. **Service Started**: Start the lender API service
4. **Test User**: Have at least two test user accounts (one owner, one non-owner)

## Start Service

```bash
# Start lender API service
cd app/lender/api
go run main.go

# Or start with specific config
go run main.go -f etc/lender.yaml
```

Service should start on port specified in config (default: 8888)

## API Endpoints

### Base URL
```
http://localhost:8888/api/v1
```

### Endpoints
1. **GET /organization/profile** - Get current user's organization info
2. **PUT /organization/profile** - Update organization information
3. **POST /organization** - Create new organization

## Authentication Setup

### Get JWT Token

First, authenticate to get a JWT token:

```bash
# Login endpoint (adjust based on your auth implementation)
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "your_password"
  }'
```

Save the `access_token` from response for subsequent requests.

### Using Token in Requests

```bash
# Set token as environment variable
export TOKEN="your_jwt_token_here"

# Use in requests
curl -H "Authorization: Bearer $TOKEN" http://localhost:8888/api/v1/organization/profile
```

## Test Cases

### T018: Test Profile API with Valid JWT Token ✅

**Objective**: Verify user can view their organization information

**Request**:
```bash
curl -X GET http://localhost:8888/api/v1/organization/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response** (200 OK):
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "organization_uuid",
    "name": "Test Organization",
    "logo": "https://example.com/logo.png",
    "website": "https://example.com",
    "description": "Test organization description",
    "status": "inactive",
    "type": "lender",
    "owner_type": "third_party",
    "owner_id": "user_uuid",
    "metadata": {
      "lending_license": "LICENSE123",
      "contact_email": "contact@example.com"
    },
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

**Verification Points**:
- [ ] Response status is 200
- [ ] Organization ID is valid UUID
- [ ] All organization fields are present
- [ ] Time fields are in RFC3339 format
- [ ] Metadata is properly serialized

---

### T019: Test UpdateProfile API with Owner Account ✅

**Objective**: Verify owner can update organization information

**Prerequisites**: Use account that is organization owner

**Request**:
```bash
curl -X PUT http://localhost:8888/api/v1/organization/profile \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Organization Name",
    "description": "Updated description",
    "website": "https://updated-site.com",
    "logo": "https://updated-site.com/new-logo.png",
    "lender_metadata": {
      "lending_license": "NEW_LICENSE456",
      "contact_email": "newcontact@example.com",
      "additional_field": "custom value"
    }
  }'
```

**Expected Response** (200 OK):
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "organization_uuid",
    "name": "Updated Organization Name",
    "description": "Updated description",
    "website": "https://updated-site.com",
    "logo": "https://updated-site.com/new-logo.png",
    "status": "inactive",
    "type": "lender",
    "owner_type": "third_party",
    "owner_id": "owner_user_uuid",
    "metadata": {
      "lending_license": "NEW_LICENSE456",
      "contact_email": "newcontact@example.com",
      "additional_field": "custom value"
    },
    "updated_at": "2025-01-09T11:00:00Z"
  }
}
```

**Verification Points**:
- [ ] Response status is 200
- [ ] Fields are updated with new values
- [ ] Unchanged fields retain original values
- [ ] Metadata is properly updated
- [ ] Updated_at timestamp changed

---

### T020: Test UpdateProfile API with Non-Owner Account ✅

**Objective**: Verify non-owner users cannot update organization

**Prerequisites**: Use account that is NOT the organization owner

**Request**:
```bash
curl -X PUT http://localhost:8888/api/v1/organization/profile \
  -H "Authorization: Bearer $NON_OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Attempted Update"
  }'
```

**Expected Response** (400/403 Error):
```json
{
  "code": 400,
  "msg": "Only Owner Can Update Organization Information"
}
```

**Verification Points**:
- [ ] Response status is 400 or 403
- [ ] Error message indicates permission denied
- [ ] Organization data is NOT updated
- [ ] Error message is clear and actionable

---

### T021: Test Create API with Valid Data ✅

**Objective**: Verify user can create new organization

**Request**:
```bash
curl -X POST http://localhost:8888/api/v1/organization \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Test Organization",
    "description": "A newly created organization",
    "website": "https://new-org.com",
    "logo": "https://new-org.com/logo.png",
    "lender_metadata": {
      "lending_license": "TEMP_LICENSE789",
      "contact_email": "admin@new-org.com"
    }
  }'
```

**Expected Response** (200 OK):
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "new_organization_uuid",
    "name": "New Test Organization",
    "description": "A newly created organization",
    "website": "https://new-org.com",
    "logo": "https://new-org.com/logo.png",
    "status": "inactive",
    "type": "lender",
    "owner_type": "third_party",
    "owner_id": "creator_user_uuid",
    "metadata": {
      "lending_license": "TEMP_LICENSE789",
      "contact_email": "admin@new-org.com"
    },
    "created_at": "2025-01-09T12:00:00Z",
    "updated_at": "2025-01-09T12:00:00Z"
  }
}
```

**Verification Points**:
- [ ] Response status is 200
- [ ] New organization ID is generated
- [ ] Creator is set as owner_id
- [ ] All provided fields are correctly saved

---

### T022: Test Create API with Missing Required Fields ✅

**Objective**: Verify validation catches missing required fields

**Request** (missing required "name" field):
```bash
curl -X POST http://localhost:8888/api/v1/organization \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Organization without name"
  }'
```

**Expected Response** (400 Error):
```json
{
  "code": 400,
  "msg": "Invalid input parameters"
}
```

**Verification Points**:
- [ ] Response status is 400
- [ ] Error message indicates validation failure
- [ ] No organization is created
- [ ] Error specifies missing required fields

---

### T023: Verify Organization Status is Inactive ✅

**Objective**: Confirm newly created organizations have status "inactive"

**Test Steps**:
1. Create a new organization (see T021)
2. Check response data.status field
3. Verify in database: `SELECT status FROM organizations WHERE id = 'new_org_uuid'`

**Expected Result**:
```json
{
  "status": "inactive"
}
```

**Database Verification**:
```sql
SELECT id, name, status, owner_type, type
FROM organizations
WHERE id = 'new_organization_uuid';
```

**Verification Points**:
- [ ] API response shows status = "inactive"
- [ ] Database record has status = "inactive"
- [ ] Status cannot be overridden in create request

---

### T024: Verify Organization Owner Type is ThirdParty ✅

**Objective**: Confirm newly created organizations have owner_type "third_party"

**Test Steps**:
1. Create a new organization (see T021)
2. Check response data.owner_type field
3. Verify in database

**Expected Result**:
```json
{
  "owner_type": "third_party"
}
```

**Database Verification**:
```sql
SELECT owner_type FROM organizations WHERE id = 'new_organization_uuid';
```

**Verification Points**:
- [ ] API response shows owner_type = "third_party"
- [ ] Database record has owner_type = "third_party"
- [ ] Owner type cannot be overridden in create request

---

### T025: Verify Creator Becomes Organization Owner ✅

**Objective**: Confirm organization creator is automatically set as owner

**Test Steps**:
1. Create a new organization with user A's token
2. Check response data.owner_id field
3. Verify owner_id matches user A's user_id
4. Check organization_users table for association

**Expected Result**:
```json
{
  "owner_id": "user_a_uuid"
}
```

**Database Verification**:
```sql
-- Check organization owner
SELECT owner_id FROM organizations WHERE id = 'new_organization_uuid';

-- Check organization user association
SELECT * FROM organization_users
WHERE organization_id = 'new_organization_uuid'
  AND user_id = 'user_a_uuid';

-- Check user has admin role
SELECT r.name
FROM organization_users ou
JOIN roles r ON ou.organization_id = r.organization_id
WHERE ou.organization_id = 'new_organization_uuid'
  AND ou.user_id = 'user_a_uuid';
```

**Verification Points**:
- [ ] owner_id matches creator's user_id
- [ ] Organization user record exists
- [ ] User has admin role in organization
- [ ] Organization user status is "active"

---

## Test Execution Checklist

Use this checklist to track test execution:

- [ ] T018: Profile API with valid JWT
- [ ] T019: UpdateProfile with owner account
- [ ] T020: UpdateProfile with non-owner account
- [ ] T021: Create API with valid data
- [ ] T022: Create API with missing fields
- [ ] T023: Verify status = inactive
- [ ] T024: Verify owner_type = third_party
- [ ] T025: Verify creator becomes owner

## Common Issues & Troubleshooting

### Issue: "User Not Belong To Any Organization"
**Cause**: User's JWT token doesn't have organization_id claim
**Solution**: Ensure user is added to an organization before testing Profile API

### Issue: "Only Owner Can Update Organization Information"
**Cause**: Non-owner user attempting to update
**Solution**: This is expected behavior for T020 test

### Issue: Database connection errors
**Cause**: Database not running or incorrect connection config
**Solution**: Check database status and `app/lender/etc/lender.yaml` configuration

### Issue: JWT token expired
**Cause**: Token has expired
**Solution**: Login again to get fresh token

## Test Data Cleanup

After testing, clean up test data:

```sql
-- Delete test organizations
DELETE FROM organization_users WHERE organization_id IN ('test_org_uuid_1', 'test_org_uuid_2');
DELETE FROM organizations WHERE id IN ('test_org_uuid_1', 'test_org_uuid_2');
```

## Success Criteria

All tests pass when:
1. ✅ Valid requests return 200 status with correct data
2. ✅ Invalid requests return appropriate error codes and messages
3. ✅ Permission checks work correctly
4. ✅ Database state matches API responses
5. ✅ Default values are correctly set

## Next Steps After Testing

If all tests pass:
1. Feature is ready for integration testing
2. Can proceed to deployment
3. Document any edge cases found

If tests fail:
1. Record failure details
2. Check logs: `tail -f app/lender/api/logs/lender.log`
3. Verify business logic implementation
4. Report issues for fixing
