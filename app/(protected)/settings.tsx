import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View, TextInput, TouchableOpacity, Alert } from "react-native";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export type RootStackParamList = {
  Dashboard: undefined;
  Settings: undefined;
  TenantManagement: { propertyId: string };
};

type SettingsNavProp = NativeStackNavigationProp<RootStackParamList, "Settings">;

export default function SettingsScreen() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [properties, setProperties] = useState<any[]>([]);
  const navigation = useNavigation();
  const navigation2 = useNavigation<SettingsNavProp>();

  useEffect(() => {
    fetchUser();
    fetchProperties();
  }, []);

  const fetchUser = async () => {
    try {
      const user = await apiRequest("/users", "GET");
      if (user) {
        setEmail(user.email);
      }
    } catch (err) {
      console.error("User fetch error:", err);
    }
  };

  const fetchProperties = async () => {
    try {
      const res = await apiRequest("/properties", "GET");
      setProperties(res || []);
    } catch (err) {
      console.error("Property fetch error:", err);
    }
  };

  const updateUser = async () => {
    try {
      await apiRequest("/users", "PUT", { email, password });
      Alert.alert("Success", "Account updated!");
    } catch (err) {
      Alert.alert("Error", "Failed to update account");
    }
  };

  const addProperty = async () => {
    Alert.prompt("Add Property", "Enter property name", async (name) => {
      if (!name) return;
      try {
        await apiRequest("/properties", "POST", { name });
        fetchProperties();
      } catch (err) {
        Alert.alert("Error", "Failed to add property");
      }
    });
  };

  const editProperty = async (id: string, currentName: string) => {
    Alert.prompt("Edit Property", "Update property name", async (name) => {
      if (!name) return;
      try {
        await apiRequest(`/properties/${id}`, "PUT", { name });
        fetchProperties();
      } catch (err) {
        Alert.alert("Error", "Failed to edit property");
      }
    }, undefined, currentName);
  };

  const deleteProperty = async (id: string) => {
    Alert.alert("Delete Property", "Are you sure?", [
      { text: "Cancel", style: "cancel" },
      {
        text: "Delete", style: "destructive", onPress: async () => {
          try {
            await apiRequest(`/properties/${id}`, "DELETE");
            fetchProperties();
          } catch (err) {
            Alert.alert("Error", "Failed to delete property");
          }
        }
      },
    ]);
  };

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <Text style={styles.title}>Settings</Text>

        {/* Edit Account */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Account</Text>
          <TextInput
            style={styles.input}
            placeholder="Email"
            value={email}
            onChangeText={setEmail}
          />
          <TextInput
            style={styles.input}
            placeholder="New Password"
            secureTextEntry
            value={password}
            onChangeText={setPassword}
          />
          <TouchableOpacity style={styles.button} onPress={updateUser}>
            <Text style={styles.buttonText}>Update Account</Text>
          </TouchableOpacity>
        </View>

        {/* Properties */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Properties</Text>
          {properties.length > 0 ? (
            properties.map((prop, idx) => (
              <View key={idx} style={styles.row}>
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>{prop.name}</Text>
                </View>
                <View style={styles.rowRight}>
                  <TouchableOpacity onPress={() => navigation2.navigate("TenantManagement", { propertyId: prop.id })}>
                    <Text style={styles.details}>View</Text>
                  </TouchableOpacity>
                  <TouchableOpacity onPress={() => editProperty(prop.id, prop.name)}>
                    <Text style={styles.details}>Edit</Text>
                  </TouchableOpacity>
                  <TouchableOpacity onPress={() => deleteProperty(prop.id)}>
                    <Text style={[styles.details, { color: "red" }]}>Delete</Text>
                  </TouchableOpacity>
                </View>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>No properties</Text>
          )}
          <TouchableOpacity style={styles.button} onPress={addProperty}>
            <Text style={styles.buttonText}>+ Add Property</Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </LinearGradient>
  );
}
