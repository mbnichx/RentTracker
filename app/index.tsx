/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

// app/index.tsx — tiny launcher that redirects to the login flow.
import { Redirect } from "expo-router";

export default function Index() {
  // Immediately redirect to the public login screen. This keeps the root
  // route minimal and centralizes the initial navigation decision here.
  return <Redirect href="/LoginScreen" />;
}
