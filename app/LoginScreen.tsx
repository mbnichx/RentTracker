import { useRouter } from "expo-router";
import React, { useState } from "react";
import {
  Alert,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";
import { LinearGradient } from "expo-linear-gradient";
import apiRequest from "../apis/client";
import { useAuth } from "../contexts/AuthContext";

export default function LoginScreen() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();
  const { login } = useAuth();

  async function handleLogin() {
    try {
      const res = await apiRequest("/login", "POST", { email, password });

      if (res.token) {
        await login(res.token);
        router.replace("/(protected)/dashboard");
      } else {
        Alert.alert("Login failed", "No token returned from server.");
      }
    } catch (err: any) {
      console.error("Login error:", err);
      Alert.alert("Error", err.message || "Login failed");
    }
  }

  return (
    <LinearGradient
      colors={["#432c83B3", "#2575fcB3"]}
      style={styles.gradient}
      start={{ x: 0, y: 0 }}
      end={{ x: 1, y: 1 }}
    >
      <KeyboardAvoidingView
        style={styles.container}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <View style={styles.inner}>
          <Text style={styles.title}>RentTracker</Text>
          <Text style={styles.subtitle}>Login to continue</Text>

          <TextInput
            placeholder="Email"
            value={email}
            onChangeText={setEmail}
            style={styles.input}
            autoCapitalize="none"
            keyboardType="email-address"
            placeholderTextColor="#999"
          />
          <TextInput
            placeholder="Password"
            value={password}
            onChangeText={setPassword}
            style={styles.input}
            secureTextEntry
            placeholderTextColor="#999"
          />

          <TouchableOpacity style={styles.button} onPress={handleLogin}>
            <Text style={styles.buttonText}>Sign In</Text>
          </TouchableOpacity>

          <Text style={styles.footerText}>
            Don’t have an account?{" "}
            <Text
              style={styles.linkText}
              onPress={() => router.push("/RegisterScreen")}
            >
              Register
            </Text>
          </Text>
        </View>
      </KeyboardAvoidingView>
    </LinearGradient>
  );
}

const styles = StyleSheet.create({
  gradient: {
    flex: 1,
  },
  container: {
    flex: 1,
    justifyContent: "center",
    padding: 20,
  },
  inner: {
    paddingHorizontal: 24,
    alignItems: "center",
  },
  title: {
    fontSize: 28,
    fontWeight: "700",
    color: "#fff",
    marginBottom: 6,
  },
  subtitle: {
    fontSize: 16,
    color: "#e0e0e0",
    marginBottom: 24,
  },
  input: {
    width: "100%",
    backgroundColor: "#ffffffee",
    borderRadius: 12,
    padding: 14,
    marginBottom: 12,
    fontSize: 16,
    color: "#333",
  },
  button: {
    width: "100%",
    backgroundColor: "#432c83B3",
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: "center",
    marginTop: 8,
  },
  buttonText: {
    color: "#ffffffff",
    fontSize: 18,
    fontWeight: "700",
  },
  footerText: {
    marginTop: 20,
    fontSize: 14,
    color: "#fff",
  },
  linkText: {
    color: "#fff",
    fontWeight: "600",
    textDecorationLine: "underline",
  },
});