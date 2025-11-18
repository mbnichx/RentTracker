import apiRequest from "./client";

// User shape used when creating or updating users via the API.
// Note: password is hashed into backend database.
type User = {
  userFirstName: string;
  userLastName: string;
  userEmail: string;
  userPhoneNumber: string;
  userPassword: string;
  userRole: string;
};

/**
 * Create a new user in the system.
 * @param user - User payload including password and role
 * @returns Promise resolving with the created user or API response
 */
export async function createUser(user: User) {
  return apiRequest("/users", "POST", user);
}

/**
 * Fetch all users. Backend may restrict results based on the calling user's
 * permissions.
 * @returns Promise resolving with an array of users
 */
export async function getUsers() {
  return apiRequest("/users/", "GET");
}

/**
 * Fetch the currently authenticated user.
 * Useful for profile screens and determining current user's permissions.
 */
export async function getCurrentUser() {
  return apiRequest("/users/me", "GET");
}

/**
 * Update an existing user. Accepts a partial `User` so callers can update a
 * subset of fields. Include an identifier field (e.g., email or an id) per
 * backend contract to select which user to update.
 */
export async function updateUser(user: Partial<User>) {
  return apiRequest("/users/update", "PUT", user);
}

/**
 * Delete a user by numeric id.
 * @param userId - Id of the user to remove
 */
export async function deleteUser(userId: number) {
  return apiRequest(`/users/delete/${userId}`, "DELETE");
}
