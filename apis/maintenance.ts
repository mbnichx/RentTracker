import apiRequest from "./client";

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

export async function createMaintenanceRequest(maintenanceRequest:MaintenanceRequest) {
  return apiRequest("/maintenanceRequests", "POST", maintenanceRequest);
}

export async function getMaintenanceRequests() {
  return apiRequest("/maintenanceRequests/", "GET");
}

export async function updateMaintenanceRequest(maintenanceRequest:MaintenanceRequest) {
  return apiRequest("/maintenanceRequests/update", "PUT", maintenanceRequest);
}

export async function deleteMaintenanceRequest(maintenanceRequestId:number) {
  return apiRequest(`/maintenanceRequests/delete/${maintenanceRequestId}`, "DELETE");
}
