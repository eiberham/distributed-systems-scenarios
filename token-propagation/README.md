## Token Propagation

#### Keycloak configuration

1. Realm:

Create realm → token-propagation

2. Client:

Clients → Create → backend-go
Client authentication: ON
Service accounts roles: ON
Direct access grants: ON
Save → tab Credentials → copy secret

3. Service Account Roles:

Tab "Service Account Roles"
Assign role → Filter by clients → realm-management
Asignar: view-users, manage-users

4. User:

Users → Create → test-user
Tab Credentials → Set password → 123456 → Temporary: OFF