import apiRequest from "./client";

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

export async function createUnit(unit:Unit) {
  return apiRequest("/units", "POST", unit);
}

export async function getUnits() {
  return apiRequest("/units/", "GET");
}

export async function updateUnit(unit:Unit) {
  return apiRequest("/units/update", "PUT", unit);
}

export async function deleteUnit(unitId:number) {
  return apiRequest(`/units/delete/${unitId}`, "DELETE");
}
