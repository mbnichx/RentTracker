
-- PROPERTIES
INSERT INTO properties (ownerUserId, propertyName, propertyStreetAddress, propertyCity, propertyState, propertyZip, propertyType, propertyYearBuilt, propertyNotes)
VALUES
(1, 'Sunny Apartments', '123 Main St', 'Springfield', 'IL', '62701', 'multi-family', 1995, 'Near park'),
(2, 'Maple Villas', '456 Oak Ave', 'Greenville', 'TX', '75401', 'apartment', 2005, 'Quiet neighborhood');

-- UNITS
INSERT INTO propertyUnits (propertyId, propertyUnitNumber, propertyUnitBeds, propertyUnitBaths, propertyUnitSqFt, propertyUnitRentDefault, propertyUnitNotes)
VALUES
(1, '1A', 2, 1, 800, 1200, 'First floor, sunny'),
(1, '1B', 3, 2, 950, 1400, 'Second floor, balcony'),
(2, '2A', 1, 1, 600, 900, 'Near pool');

-- TENANTS
INSERT INTO tenants (tenantFirstName, tenantLastName, tenantEmailAddress, tenantPhoneNumber)
VALUES
('John', 'Doe', 'john@example.com', '555-3000'),
('Jane', 'Smith', 'jane@example.com', '555-4000'),
('Sam', 'Brown', 'sam@example.com', '555-5000');

-- LEASES (UNIX timestamps: example current time ~ 1700000000)
INSERT INTO leases (tenantId, propertyUnitId, leaseStartUnix, leaseEndUnix, leaseRentAmount, leaseSecurityDeposit, leaseDocumentLink, leaseStatus)
VALUES
(1, 1, 1690000000, 1720000000, 1200, 1200, 'link1', 'active'),
(2, 2, 1695000000, 1725000000, 1400, 1400, 'link2', 'active'),
(3, 3, 1698000000, 1728000000, 900, 900, 'link3', 'active');

-- PAYMENTS
INSERT INTO payments (leaseId, paymentAmount, paymentDateUnix, paymentMethod, paymentNotes)
VALUES
(1, 1200, 1690300000, 'ACH', 'Paid on time'),
(1, 1200, 1693000000, 'ACH', 'Overdue'),
(2, 1400, 1696000000, 'Check', 'Early payment'),
(3, 900, 1699000000, 'Cash', 'Paid');

-- MAINTENANCE REQUESTS
INSERT INTO maintenanceRequests (propertyUnitId, leaseId, maintenanceRequestInfo, maintenanceRequestPriority, maintenanceRequestCategory, maintenanceRequestStatus, maintenanceRequestCreatedUnix, maintenanceRequestCompletedUnix, maintenanceAssignedTo)
VALUES
(1, 1, 'Leaky faucet in kitchen', 'normal', 'plumbing', 'completed', 1690500000, 1690600000, 'Bob'),
(2, 2, 'AC not cooling', 'urgent', 'HVAC', 'open', 1696100000, NULL, 'Alice'),
(3, 3, 'Broken window', 'low', 'general', 'in progress', 1699100000, NULL, 'Bob');