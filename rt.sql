-- USERS (landlords, managers, assistants, etc.)
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    userId INTEGER PRIMARY KEY AUTOINCREMENT,
    userFirstName TEXT NOT NULL,
    userLastName TEXT NOT NULL,
    userEmail TEXT UNIQUE NOT NULL,
    userPhoneNumber TEXT,
    userPasswordHash TEXT NOT NULL,
    userRole TEXT DEFAULT 'owner' -- owner, manager, assistant
);

-- PROPERTIES (owned by users)
DROP TABLE IF EXISTS properties;
CREATE TABLE IF NOT EXISTS properties (
    propertyId INTEGER PRIMARY KEY AUTOINCREMENT,
    ownerUserId INTEGER REFERENCES users(userId),
    propertyName TEXT,
    propertyStreetAddress TEXT,
    propertyCity TEXT,
    propertyState TEXT,
    propertyZip TEXT,
    propertyType TEXT, -- single-family, multi-family, apartment
    propertyYearBuilt INTEGER,
    propertyNotes TEXT
);

-- UNITS (belonging to properties)
DROP TABLE IF EXISTS propertyUnits;
CREATE TABLE IF NOT EXISTS propertyUnits (
    propertyUnitId INTEGER PRIMARY KEY AUTOINCREMENT,
    propertyId INTEGER REFERENCES properties(propertyId),
    propertyUnitNumber TEXT,
    propertyUnitBeds INTEGER,
    propertyUnitBaths INTEGER,
    propertyUnitSqFt INTEGER,
    propertyUnitRentDefault INTEGER,
    propertyUnitNotes TEXT
);

-- TENANTS (people who sign leases)
DROP TABLE IF EXISTS tenants;
CREATE TABLE IF NOT EXISTS tenants (
    tenantId INTEGER PRIMARY KEY AUTOINCREMENT,
    tenantFirstName TEXT NOT NULL,
    tenantLastName TEXT NOT NULL,
    tenantEmailAddress TEXT,
    tenantPhoneNumber TEXT
);

-- LEASES (link tenants to units over time)
DROP TABLE IF EXISTS leases;
CREATE TABLE IF NOT EXISTS leases (
    leaseId INTEGER PRIMARY KEY AUTOINCREMENT,
    tenantId INTEGER REFERENCES tenants(tenantId),
    propertyUnitId INTEGER REFERENCES propertyUnits(propertyUnitId),
    leaseStartUnix INTEGER NOT NULL,
    leaseEndUnix INTEGER,
    leaseRentAmount INTEGER NOT NULL,
    leaseSecurityDeposit INTEGER,
    leaseDocumentLink TEXT,
    leaseStatus TEXT DEFAULT 'active' -- active, expired, renewed
);

-- PAYMENTS (linked to leases, not directly to tenants)
DROP TABLE IF EXISTS payments;
CREATE TABLE IF NOT EXISTS payments (
    paymentId INTEGER PRIMARY KEY AUTOINCREMENT,
    leaseId INTEGER REFERENCES leases(leaseId),
    paymentAmount INTEGER NOT NULL,
    paymentDateUnix INTEGER NOT NULL,
    paymentMethod TEXT, -- cash, check, ACH, Zelle, etc.
    paymentNotes TEXT,
    paymentConfirmation BLOB
);

-- MAINTENANCE REQUESTS (per unit, may involve lease but usually tied to unit)
DROP TABLE IF EXISTS maintenanceRequests;
CREATE TABLE IF NOT EXISTS maintenanceRequests (
    maintenanceRequestId INTEGER PRIMARY KEY AUTOINCREMENT,
    propertyUnitId INTEGER REFERENCES propertyUnits(propertyUnitId),
    leaseId INTEGER REFERENCES leases(leaseId), -- optional, if tenant reported
    maintenanceRequestInfo TEXT NOT NULL,
    maintenanceRequestPriority TEXT DEFAULT 'normal', -- low, normal, urgent
    maintenanceRequestCategory TEXT, -- plumbing, HVAC, etc.
    maintenanceRequestStatus TEXT DEFAULT 'open', -- open, in progress, completed
    maintenanceRequestCreatedUnix INTEGER NOT NULL,
    maintenanceRequestCompletedUnix INTEGER,
    maintenanceAssignedTo TEXT
);

-- ACTIVITY LOG (audit trail, optional)
DROP TABLE IF EXISTS activityLogs;
CREATE TABLE IF NOT EXISTS activityLogs (
    logId INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER REFERENCES users(userId),
    entityType TEXT, -- e.g. 'payment', 'maintenance', 'lease'
    entityId INTEGER,
    action TEXT, -- created, updated, deleted
    timestampUnix INTEGER NOT NULL
);
