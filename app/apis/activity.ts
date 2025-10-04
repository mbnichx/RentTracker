import apiRequest from "./client";

type Activity = {
  logId: number;
  userId: number;
  entityType: string;
  entityId: number;
  action: string;
  timeStampUnix: number;
};

export async function createActivity(activity:Activity) {
  return apiRequest("/activities", "POST", activity);
}

export async function getActivitys() {
  return apiRequest("/activities/", "GET");
}

export async function updateActivity(activity:Activity) {
  return apiRequest("/activities/update", "PUT", activity);
}

export async function deleteActivity(activityId:number) {
  return apiRequest(`/activities/delete/${activityId}`, "DELETE");
}
