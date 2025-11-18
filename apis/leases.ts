/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

import apiRequest from "./client";

// Shape of a lease record exchanged with the backend.
type Lease = {
  leaseId: number;
  tenantId: number;
  propertyUnitId: number;
  leaseStartUnix: number;
  leaseEndUnix: number;
  leaseRentAmount: number;
  leaseSecurityDeposit: number;
  leaseDocumentLink: string;
  leaseStatus: string;
};

/**
 * Create a lease record on the server.
 * @param lease - Lease object to create
 * @returns Promise resolving with server response
 */
export async function createLease(lease: Lease) {
  return apiRequest("/leases", "POST", lease);
}

/**
 * Retrieve all leases.
 * @returns Promise resolving with an array of leases
 */
export async function getLeases() {
  return apiRequest("/leases/", "GET");
}

/**
 * Update an existing lease. The `leaseId` should be present on the payload
 * to identify which record to update.
 * @param lease - Updated lease object
 */
export async function updateLease(lease: Lease) {
  return apiRequest("/leases/update", "PUT", lease);
}

/**
 * Delete a lease by id.
 * @param leaseId - Numeric id of the lease to delete
 */
export async function deleteLease(leaseId: number) {
  return apiRequest(`/leases/delete/${leaseId}`, "DELETE");
}
