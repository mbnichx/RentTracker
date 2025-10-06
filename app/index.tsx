// app/index.tsx
import { Redirect } from "expo-router";

export default function Index() {
  // When app launches, automatically redirect to LoginScreen
  return <Redirect href="/LoginScreen" />;
}
