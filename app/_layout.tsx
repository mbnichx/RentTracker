/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Root application layout.
 *
 * This file wraps the application's Router `Stack` with the `AuthProvider`
 * so authentication state is available throughout the app. The provider
 * handles token storage and exposes hooks (see `contexts/AuthContext`).
 */
import { Stack } from "expo-router";
import { AuthProvider } from "../contexts/AuthContext";

export default function RootLayout() {
  // The `Stack` component renders the router-controlled screens. We hide
  // native headers because each screen uses custom header UI or no header.
  return (
    <AuthProvider>
      <Stack screenOptions={{ headerShown: false }} />
    </AuthProvider>
  );
}
