import apiRequest from "./client";

export type Property = {
  propertyId: number;
  ownerUserId: number;
  propertyName: string;
  propertyStreetAddress: string;
  propertyCity: string;
  propertyState: string;
  propertyZip: string;
  propertyType: string;
  propertyYearBuilt: string;
  propertyNotes: string;
};

export async function createProperty(property:Partial<Property>) {
  return apiRequest("/properties", "POST", property);
}

export async function getProperties() {
  return apiRequest("/properties/", "GET");
}

export async function updateProperty(property:Partial<Property>) {
  return apiRequest("/properties/update", "PUT", property);
}

export async function deleteProperty(propertyId:number) {
  return apiRequest(`/properties/delete/${propertyId}`, "DELETE");
}
