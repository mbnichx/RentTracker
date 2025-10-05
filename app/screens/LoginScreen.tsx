import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert } from "react-native";
import { useRouter } from "expo-router";

const API_URL = "http://10.0.0.67:8080";

export default function LoginScreen() {
  const router = useRouter();
  const [userEmail, setUserEmail] = useState("");
  const [userPassword, setUserPassword] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleLogin() {
    if (!userEmail || !userPassword) {
      Alert.alert("Error", "Please enter both email and password.");
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userEmail, userPassword }),
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || "Login failed");
      }

      const data = await res.json();
      console.log("✅ Login successful:", data);
      Alert.alert("Success", "Logged in successfully!");
      router.replace("/screens/settings"); // Navigate to home screen (create app/home.tsx if not yet)
    } catch (err: any) {
      console.error("Login error:", err);
      Alert.alert("Error", err.message || "Could not connect to server");
    } finally {
      setLoading(false);
    }
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Landlord Login</Text>

      <TextInput
        style={styles.input}
        placeholder="Email"
        value={userEmail}
        onChangeText={setUserEmail}
        autoCapitalize="none"
      />

      <TextInput
        style={styles.input}
        placeholder="Password"
        secureTextEntry
        value={userPassword}
        onChangeText={setUserPassword}
      />

      <TouchableOpacity style={styles.button} onPress={handleLogin} disabled={loading}>
        <Text style={styles.buttonText}>{loading ? "Logging in..." : "Login"}</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={() => router.push("./RegisterScreen" as never)}>
        <Text style={styles.link}>Don’t have an account? Register</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: "#fff", padding: 20 },
  title: { fontSize: 26, fontWeight: "bold", marginBottom: 24 },
  input: {
    width: "100%",
    borderWidth: 1,
    borderColor: "#ccc",
    borderRadius: 8,
    padding: 12,
    marginBottom: 16,
  },
  button: { backgroundColor: "#007AFF", padding: 14, borderRadius: 8, width: "100%", alignItems: "center" },
  buttonText: { color: "#fff", fontWeight: "bold", fontSize: 16 },
  link: { marginTop: 16, color: "#007AFF", fontWeight: "500" },
});
