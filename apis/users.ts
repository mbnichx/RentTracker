import apiRequest from "./client";

type User = {
  userFirstName: string;
  userLastName: string;
  userEmail: string;
  userPhoneNumber: string;
  userPassword: string;
  userRole: string;
};

export async function createUser(user:User) {
  console.log(user)
  return apiRequest("/users", "POST", user);
}

export async function getUsers() {
  return apiRequest("/users/", "GET");
}

export async function updateUser(user:User) {
  return apiRequest("/users/update", "PUT", user);
}

export async function deleteUser(userId:number) {
  return apiRequest(`/users/delete/${userId}`, "DELETE");
}
