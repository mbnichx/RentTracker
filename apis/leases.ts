import apiRequest from "./client";

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

export async function createLease(lease:Lease) {
  return apiRequest("/leases", "POST", lease);
}

export async function getLeases() {
  return apiRequest("/leases/", "GET");
}

export async function updateLease(lease:Lease) {
  return apiRequest("/leases/update", "PUT", lease);
}

export async function deleteLease(leaseId:number) {
  return apiRequest(`/leases/delete/${leaseId}`, "DELETE");
}
