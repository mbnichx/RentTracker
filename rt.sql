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


-- == Views =====================================================================
-- Overdue Rent (dashboard)
DROP VIEW IF EXISTS overduePayments;

CREATE VIEW overduePayments AS
SELECT 
    t.tenantFirstName AS firstName,
    t.tenantLastName AS lastName,
    p.propertyStreetAddress AS address,
    u.propertyUnitNumber AS unit,
    l.leaseRentAmount AS rentAmount,
    COALESCE(MAX(pay.paymentDateUnix), 0) AS lastPaymentUnix,
    CASE 
        WHEN MAX(pay.paymentDateUnix) IS NULL THEN 'No payment recorded'
        WHEN (strftime('%s','now') - MAX(pay.paymentDateUnix)) > (30 * 24 * 60 * 60) THEN 'Overdue'
        ELSE 'Current'
    END AS paymentStatus
FROM leases l
JOIN tenants t ON l.tenantId = t.tenantId
JOIN propertyUnits u ON l.propertyUnitId = u.propertyUnitId
JOIN properties p ON u.propertyId = p.propertyId
LEFT JOIN payments pay ON l.leaseId = pay.leaseId
WHERE l.leaseStatus = 'active'
GROUP BY l.leaseId
HAVING paymentStatus IN ('Overdue', 'No payment recorded');


-- Maintenance requests (dashboard)
DROP VIEW IF EXISTS maintenanceRequestsView;
CREATE VIEW maintenanceRequestsView AS
SELECT
    t.tenantFirstName AS firstName,
    t.tenantLastName AS lastName,
    p.propertyStreetAddress AS address,
    u.propertyUnitNumber AS unit,
    m.maintenanceRequestInfo AS description,
    m.maintenanceRequestStatus AS maintenanceStatus,
    m.maintenanceRequestCreatedUnix AS dateCreated,
    m.maintenanceRequestPriority AS priority,
    m.maintenanceRequestCategory AS category
FROM maintenanceRequests m
LEFT JOIN leases l ON m.leaseId = l.leaseId
LEFT JOIN tenants t ON l.tenantId = t.tenantId
LEFT JOIN propertyUnits u ON m.propertyUnitId = u.propertyUnitId
LEFT JOIN properties p ON u.propertyId = p.propertyId
WHERE m.maintenanceRequestStatus != 'completed';

-- Lease renewals (dashboard)
DROP VIEW IF EXISTS leasesView;
CREATE VIEW leasesView AS
SELECT
    t.tenantFirstName AS firstName,
    t.tenantLastName AS lastName,
    p.propertyStreetAddress AS address,
    u.propertyUnitNumber AS unit,
    l.leaseStartUnix AS leaseStartDate,
    l.leaseRentAmount AS rentAmount,
    l.leaseStatus AS leaseStatus
FROM leases l
JOIN tenants t ON l.tenantId = t.tenantId
JOIN propertyUnits u ON l.propertyUnitId = u.propertyUnitId
JOIN properties p ON u.propertyId = p.propertyId;

-- Upcoming Rent (next due payments)
DROP VIEW IF EXISTS upcomingPayments;

CREATE VIEW upcomingPayments AS
SELECT 
    t.tenantFirstName AS firstName,
    t.tenantLastName AS lastName,
    p.propertyStreetAddress AS address,
    u.propertyUnitNumber AS unit,
    l.leaseRentAmount AS rentAmount,
    COALESCE(MAX(pay.paymentDateUnix), 0) AS lastPaymentUnix,
    CASE 
        WHEN MAX(pay.paymentDateUnix) >= strftime('%s','now') THEN 'Paid'
        ELSE 'Due'
    END AS paymentStatus
FROM leases l
JOIN tenants t ON l.tenantId = t.tenantId
JOIN propertyUnits u ON l.propertyUnitId = u.propertyUnitId
JOIN properties p ON u.propertyId = p.propertyId
LEFT JOIN payments pay ON l.leaseId = pay.leaseId
WHERE l.leaseStatus = 'active'
GROUP BY l.leaseId;
