import apiRequest from "./client";

type Property = {
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

export async function createProperty(property:Property) {
  return apiRequest("/properties", "POST", property);
}

export async function getPropertys() {
  return apiRequest("/properties/", "GET");
}

export async function updateProperty(property:Property) {
  return apiRequest("/properties/update", "PUT", property);
}

export async function deleteProperty(propertyId:number) {
  return apiRequest(`/properties/delete/${propertyId}`, "DELETE");
}
