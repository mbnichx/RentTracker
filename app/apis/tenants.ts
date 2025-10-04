import apiRequest from "./client";

type Tenant = {
  tenantId: number;
  tenantFirstName: string;
  tenantLastName: string;
  tenantEmailAddress: string;
  tenantPhoneNumber: string;
};

export async function createTenant(tenant:Tenant) {
  return apiRequest("/tenants", "POST", tenant);
}

export async function getTenants() {
  return apiRequest("/tenants/", "GET");
}

export async function updateTenant(tenant:Tenant) {
  return apiRequest("/tenants/update", "PUT", tenant);
}

export async function deleteTenant(tenantId:number) {
  return apiRequest(`/tenants/delete/${tenantId}`, "DELETE");
}
