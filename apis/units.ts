/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

import apiRequest from "./client";

// Unit represents a rentable unit within a property.
type Unit = {
  propertyUnitId: number;
  propertyId: number;
  propertyUnitNumber: string;
  propertyUnitBeds: number;
  propertyUnitBaths: number;
  propertySqFt: number;
  propertyRentDefault: number;
  propertyUnitNotes: string;
};

/**
 * Create a unit record.
 * @param unit - Unit payload
 */
export async function createUnit(unit: Unit) {
  return apiRequest("/units", "POST", unit);
}

/**
 * Get units for a specific property.
 * @param propertyId - Id of property to fetch units for
 */
export async function getUnits(propertyId: number) {
  return apiRequest(`/units/${propertyId}`, "GET");
}

/**
 * Update a unit record. Include `propertyUnitId` in the unit payload.
 */
export async function updateUnit(unit: Unit) {
  return apiRequest("/units/update", "PUT", unit);
}

/**
 * Delete a unit by id.
 */
export async function deleteUnit(unitId: number) {
  return apiRequest(`/units/delete/${unitId}`, "DELETE");
}
