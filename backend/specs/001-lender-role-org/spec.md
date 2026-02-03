# Feature Specification: Lender Role and Organization User Management

**Feature Branch**: `001-lender-role-org`
**Created**: 2025-12-29
**Status**: ✅ Completed
**Completed**: 2025-12-29
**Input**: User description: "参照borrower服务将borrower中的role，organization
role功能在lender端也实现一遍，包括邀请组织用户管理，邀请组织用户加入组织等功能"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Role Management for Lender Organizations (Priority: P1)

As a lender organization administrator, I need to manage roles within my organization so that I can define what
different team members can and cannot do in the system.

**Why this priority**: Role management is foundational to any permission system. Without defined roles, we cannot assign
granular permissions to users, making the system insecure and unusable for teams.

**Independent Test**: Can be fully tested by creating, reading, updating, and deleting roles through the management
interface, and verifying that role changes are correctly persisted and retrievable.

**Acceptance Scenarios**:

1. **Given** I am an organization administrator, **When** I create a new role with name "Loan Officer" and assign
   specific permissions, **Then** the role is created successfully with the specified permissions
2. **Given** an existing role, **When** I update its name or associated permissions, **Then** the changes are saved and
   visible in subsequent queries
3. **Given** an active role assigned to users, **When** I attempt to delete the role, **Then** the system prevents
   deletion and informs me that the role is in use
4. **Given** a role with no assigned users, **When** I delete the role, **Then** the role is permanently removed from
   the system

---

### User Story 2 - Organization User Management (Priority: P2)

As a lender organization administrator, I need to manage users within my organization so that I can control who has
access to my organization's data and what they can do.

**Why this priority**: User management enables team collaboration. Administrators need to add team members, update their
information, and control their access status. This is essential for day-to-day operations.

**Independent Test**: Can be fully tested by adding users to an organization, viewing paginated user lists, updating
user details, and managing user active/disabled status.

**Acceptance Scenarios**:

1. **Given** I am an organization administrator, **When** I add a new user by email with assigned roles, **Then** the
   user is created in pending status and receives an invitation email
2. **Given** multiple users in my organization, **When** I search or filter users by keyword, status, or role, **Then**
   I see the correct filtered and paginated list of users
3. **Given** an existing organization user, **When** I update their name, phone, or assigned roles, **Then** the changes
   are saved immediately
4. **Given** an active organization user, **When** I disable their account, **Then** the user can no longer access the
   organization's resources

---

### User Story 3 - Invite and Accept Organization Membership (Priority: P3)

As a lender organization administrator, I need to invite new team members to join my organization so that they can
create their accounts and start working with our lending operations.

As a invited user, I need to accept the invitation and complete my profile so that I can access the organization's
system.

**Why this priority**: The invitation flow enables secure onboarding of new team members. It ensures that only
authorized individuals with valid email addresses can join an organization.

**Independent Test**: Can be fully tested by sending an invitation, receiving the invitation code/email, and using that
code to accept and create the user account.

**Acceptance Scenarios**:

1. **Given** I am an organization administrator, **When** I invite a user by email address, **Then** an invitation is
   created and the user receives an email with an invitation code
2. **Given** I received an invitation email, **When** I click the invitation link or enter the code, **Then** I am
   prompted to create my account with my name and email
3. **Given** I complete the registration with the invitation code, **When** I submit my information, **Then** my account
   is created and I am added to the organization with active status
4. **Given** an invitation that was already used, **When** someone attempts to accept it again, **Then** the system
   rejects the attempt with an appropriate error message

---

### User Story 4 - Role Assignment and Permission Updates (Priority: P4)

As a lender organization administrator, I need to assign or change roles for organization users so that team members can
have appropriate access levels as their responsibilities change.

**Why this priority**: Role assignment enables granular access control. Users may need different permissions as they
take on new responsibilities within the organization.

**Independent Test**: Can be fully tested by assigning multiple roles to a user, removing roles, and verifying that the
user's effective permissions match their assigned roles.

**Acceptance Scenarios**:

1. **Given** an organization user, **When** I assign one or more roles to them, **Then** the user's permissions are
   updated to include all permissions from assigned roles
2. **Given** a user with multiple roles, **When** I remove a role assignment, **Then** the user no longer has the
   permissions associated with that role
3. **Given** a user with no roles assigned, **When** they attempt to access protected resources, **Then** access is
   denied with an appropriate message

---

### Edge Cases

- What happens when an administrator attempts to delete a role that is currently assigned to one or more users?
- What happens when a user tries to accept an invitation with an expired or invalid code?
- What happens when an administrator attempts to assign a role that doesn't exist in the organization?
- What happens when a user is disabled while they are actively logged in?
- What happens when an organization has no roles defined and a new user is added?
- What happens when invitation emails fail to deliver (email service issues)?
- What happens when two administrators simultaneously update the same user's roles?
- What happens when a user is invited to multiple organizations with the same email address?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow organization administrators to create custom roles with unique names and descriptions
- **FR-002**: System MUST allow administrators to associate permissions/resources with roles during creation and updates
- **FR-003**: System MUST prevent creation of roles with names that conflict with system reserved names (e.g., "admin")
- **FR-004**: System MUST allow administrators to view all roles in their organization in a paginated list
- **FR-005**: System MUST allow administrators to filter and sort roles by name, status, or creation date
- **FR-006**: System MUST allow administrators to update role name, description, status, and associated permissions
- **FR-007**: System MUST prevent deletion of roles that are currently assigned to one or more users
- **FR-008**: System MUST allow deletion of roles that have no assigned users
- **FR-009**: System MUST allow administrators to activate or deactivate roles without deleting them
- **FR-010**: System MUST allow administrators to add new users to their organization by email address
- **FR-011**: System MUST generate a unique invitation code when a new user is added to an organization
- **FR-012**: System MUST send an invitation email to the specified email address with the invitation code
- **FR-013**: System MUST create the new user in pending status until they accept the invitation
- **FR-014**: System MUST allow administrators to view all users in their organization with pagination
- **FR-015**: System MUST allow filtering of organization users by keyword (name/email), status, or assigned role
- **FR-016**: System MUST allow sorting of organization users by name, email, status, or creation date
- **FR-017**: System MUST allow administrators to update user information (name, phone, remark)
- **FR-018**: System MUST allow administrators to assign or update roles for organization users
- **FR-019**: System MUST allow administrators to enable or disable user accounts
- **FR-020**: System MUST allow invited users to accept invitations by providing the invitation code
- **FR-021**: System MUST require invited users to provide their first name, last name, and email during acceptance
- **FR-022**: System MUST activate user accounts upon successful invitation acceptance
- **FR-023**: System MUST prevent reuse of invitation codes after they have been used
- **FR-024**: System MUST mark invitation codes as expired after a configurable time period
- **FR-025**: System MUST allow users to view their assigned roles and permissions
- **FR-026**: System MUST maintain an audit trail of role changes and user management actions

### Key Entities

- **Role**: Represents a collection of permissions that can be assigned to users. Attributes include name, description,
  status (Active/Inactive), associated resources/permissions, organization association, and timestamps.

- **Organization User**: Represents a user's membership in an organization. Attributes include user ID, organization ID,
  status (Active/Disabled/Pending), personal information (first name, last name, phone, email), assigned roles, and
  timestamps.

- **Permission/Resource**: Represents a specific capability or access right in the system. Resources define what actions
  can be performed (e.g., create loan, view reports, manage users).

- **Invitation**: Represents a pending invitation for a user to join an organization. Attributes include invitation
  code, email address, organization ID, associated roles, expiration time, and status.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Administrators can create a new role and assign permissions within 2 minutes
- **SC-002**: Organization user list loads within 3 seconds for organizations with up to 1000 users
- **SC-003**: Invitation emails are delivered within 30 seconds of user creation
- **SC-004**: 95% of invited users successfully complete registration on first attempt
- **SC-005**: Role permission changes take effect within 5 seconds for affected users
- **SC-006**: System prevents 100% of attempts to delete roles with active user assignments
- **SC-007**: Users can accept invitations and create accounts in under 3 minutes
- **SC-008**: Search and filter operations return results within 2 seconds
- **SC-009**: 100% of role-based permission checks are correctly enforced during user actions
- **SC-010**: Audit trail captures 100% of role and user management changes

## Assumptions

1. The lender service already has basic organization and user authentication infrastructure in place
2. An email service is already configured for sending notifications
3. The borrower service's role and organization user implementation serves as the reference architecture
4. Roles and permissions are organization-scoped (not global across all organizations)
5. Each organization has at least one administrator with full permissions
6. Invitation codes have a default expiration of 7 days unless otherwise configured
7. Users can belong to multiple organizations but have separate role assignments per organization
8. The system already has authentication middleware to identify the current user and their organization context

## Dependencies

1. **Existing lender authentication system** - must provide user identity and organization context
2. **Email notification service** - must be functional for sending invitation emails
3. **Organization management** - lender service must have existing organization entities and management
4. **Resource/permission definitions** - system must have a defined set of resources and actions that can be permitted
5. **Database schema** - must support tables for roles, user-role relationships, and invitations
