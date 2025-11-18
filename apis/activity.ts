// Thin API wrapper for activity-related endpoints.
// Uses the shared `apiRequest` helper to perform HTTP requests to the backend.
import apiRequest from "./client";

// Represents an activity log entry recorded by the backend
// Note: kept as an internal 'type' to avoid changing module exports
type Activity = {
  logId: number;
  userId: number;
  entityType: string; 
  entityId: number; 
  action: string; 
  timeStampUnix: number; 
};

/**
 * Create a new activity log entry.
 * @param activity - Activity object containing details to store
 * @returns Promise resolving with the API response
 */
export async function createActivity(activity: Activity) {
  // POST /activities with activity payload
  return apiRequest("/activities", "POST", activity);
}

/**
 * Fetch all activity entries.
 * Note: endpoint returns an array of activities (server-side shape).
 * @returns Promise resolving with the list of activities
 */
export async function getActivitys() {
  // GET /activities/
  return apiRequest("/activities/", "GET");
}

/**
 * Update an existing activity entry.
 * @param activity - Activity object (should include `logId` to identify the record)
 * @returns Promise resolving with the API response
 */
export async function updateActivity(activity: Activity) {
  // PUT /activities/update with updated activity payload
  return apiRequest("/activities/update", "PUT", activity);
}

/**
 * Delete an activity by id.
 * @param activityId - numeric id of the activity to delete
 * @returns Promise resolving with the API response
 */
export async function deleteActivity(activityId: number) {
  // DELETE /activities/delete/:id
  return apiRequest(`/activities/delete/${activityId}`, "DELETE");
}
