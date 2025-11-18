/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

import apiRequest from "./client";

// Represents a maintenance request ticket in the system.
type MaintenanceRequest = {
  maintenanceRequestId: number;
  propertyUnitId: number;
  leaseId: number;
  maintenanceRequestInfo: string;
  maintenanceRequestPriority: string;
  maintenanceRequestCategory: string;
  maintenanceRequestStatus: string;
  maintenanceRequestCreatedUnix: number;
  maintenanceRequestCompletedUnix: number;
  maintenanceAssignedTo: string;
};

/**
 * Create a new maintenance request.
 * @param maintenanceRequest - The maintenance payload to create
 */
export async function createMaintenanceRequest(maintenanceRequest: MaintenanceRequest) {
  return apiRequest("/maintenanceRequests", "POST", maintenanceRequest);
}

/**
 * Fetch all maintenance requests.
 * @returns Promise resolving with an array of maintenance requests
 */
export async function getMaintenanceRequests() {
  return apiRequest("/maintenanceRequests/", "GET");
}

/**
 * Update an existing maintenance request. The payload should contain the
 * `maintenanceRequestId` to identify the record to modify.
 */
export async function updateMaintenanceRequest(maintenanceRequest: MaintenanceRequest) {
  return apiRequest("/maintenanceRequests/update", "PUT", maintenanceRequest);
}

/**
 * Delete a maintenance request by id.
 */
export async function deleteMaintenanceRequest(maintenanceRequestId: number) {
  return apiRequest(`/maintenanceRequests/delete/${maintenanceRequestId}`, "DELETE");
}
