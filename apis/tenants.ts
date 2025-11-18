import apiRequest from "./client";

// Tenant record structure used by the API layer.
type Tenant = {
  tenantId: number;
  tenantFirstName: string;
  tenantLastName: string;
  tenantEmailAddress: string;
  tenantPhoneNumber: string;
};

/**
 * Create a tenant record.
 * @param tenant - Tenant payload
 */
export async function createTenant(tenant: Tenant) {
  return apiRequest("/tenants", "POST", tenant);
}

/**
 * Retrieve tenants.
 * @returns Promise resolving with an array of tenants
 */
export async function getTenants() {
  return apiRequest("/tenants/", "GET");
}

/**
 * Update a tenant record. Include `tenantId` in the payload to identify the
 * record to update.
 */
export async function updateTenant(tenant: Tenant) {
  return apiRequest("/tenants/update", "PUT", tenant);
}

/**
 * Delete a tenant by id.
 */
export async function deleteTenant(tenantId: number) {
  return apiRequest(`/tenants/delete/${tenantId}`, "DELETE");
}
