/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

import apiRequest from "./client";

/**
 * Public Property type exported for use across the app. Fields mirror the
 * server-side property record.
 */
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

/**
 * Create a property. The API accepts a partial object so callers only need to
 * provide fields required for creation.
 * @param property - Partial Property payload
 */
export async function createProperty(property: Partial<Property>) {
  return apiRequest("/properties", "POST", property);
}

/**
 * Retrieve properties. Depending on backend, this may return all properties
 * for the current user or the entire dataset.
 */
export async function getProperties() {
  return apiRequest("/properties/", "GET");
}

/**
 * Update a property. Provide the propertyId in the partial to identify the
 * record to modify.
 */
export async function updateProperty(property: Partial<Property>) {
  return apiRequest("/properties/update", "PUT", property);
}

/**
 * Delete a property by id.
 */
export async function deleteProperty(propertyId: number) {
  return apiRequest(`/properties/delete/${propertyId}`, "DELETE");
}
