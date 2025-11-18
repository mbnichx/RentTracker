/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Settings screen — allows the user to update their account and manage
 * properties. Uses `safeCall` from the api client to show friendly alerts on
 * success and handle errors in a centralized way.
 */
import { useFocusEffect, useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { LinearGradient } from "expo-linear-gradient";
import { router } from "expo-router";
import React, { useCallback, useEffect, useState } from "react";
import {
  Alert,
  ScrollView,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from "react-native";

import * as util from "../../apis/client";
import * as p from "../../apis/properties";
import * as u from "../../apis/users";
import { styles } from "./style";

// Navigation params used when pushing to tenant management
type RootStackParamList = {
  TenantManagement: { propertyId: number };
};

export default function SettingsScreen() {
  const navigation =
    useNavigation<NativeStackNavigationProp<RootStackParamList>>();

  const [userEmail, setEmail] = useState("");
  const [userPassword, setPassword] = useState("");
  // Use the exported Property type from the properties API for stronger typing
  const [properties, setProperties] = useState<p.Property[]>([]);
  const [newPropertyName, setNewPropertyName] = useState("");

  // Load user + properties on mount. `safeCall` shows alerts and returns the
  // API result or undefined on error, keeping the UI simple.
  const fetchSettingsData = useCallback(async () => {
    const user = await util.safeCall(() => u.getUsers());
    if (user?.email) setEmail(user.email);

    const props = await util.safeCall(() => p.getProperties());
    if (props) setProperties(props);
  }, []);

  useEffect(() => {
    fetchSettingsData();
  }, [fetchSettingsData]);

  useFocusEffect(
    useCallback(() => {
      fetchSettingsData();
    }, [fetchSettingsData])
  );


  const handleUpdateUser = () =>
    util.safeCall(() => u.updateUser({ userEmail: userEmail, userPassword: userPassword }), "Account updated");

  const handleAddProperty = () =>
    newPropertyName.trim()
      ? util.safeCall(() => p.createProperty({ propertyName: newPropertyName }), "Property added")
        .then((created) => created && setProperties((prev) => [...prev, created]))
        .then(() => setNewPropertyName(""))
      : Alert.alert("Validation", "Property name is required");

  const handleEditProperty = (propertyId: number, propertyName: string) =>
    util.safeCall(() => p.updateProperty({ propertyId, propertyName }), "Property updated")
      .then(() =>
        setProperties((prev) =>
          prev.map((prop) =>
            prop.propertyId === propertyId ? { ...prop, propertyName } : prop
          )
        )
      );

  const handleDeleteProperty = (propertyId: number) =>
    util.safeCall(() => p.deleteProperty(propertyId), "Property deleted")
      .then(() =>
        setProperties((prev) => prev.filter((pr) => pr.propertyId !== propertyId))
      );

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <Text style={styles.title}>Settings</Text>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>Account</Text>
          <TextInput
            style={styles.input}
            placeholder="Email"
            value={userEmail}
            onChangeText={setEmail}
            autoCapitalize="none"
          />
          <TextInput
            style={styles.input}
            placeholder="New Password"
            secureTextEntry
            value={userPassword}
            onChangeText={setPassword}
          />
          <TouchableOpacity style={styles.button} onPress={handleUpdateUser}>
            <Text style={styles.buttonText}>Update Account</Text>
          </TouchableOpacity>
        </View>

        {/* Properties */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Properties</Text>
          {properties.length > 0 ? (
            properties.map((prop) => (
              <View key={prop.propertyId} style={styles.row}>
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>{prop.propertyName}</Text>
                </View>
                <View style={styles.rowRight}>
                  <TouchableOpacity
                    onPress={() =>
                      router.push({
                        pathname: "/settingsTenantMgmt",
                        params: { propertyId: prop.propertyId }
                      })
                    }
                  >
                    <Text style={styles.details}>View</Text>
                  </TouchableOpacity>
                  <TouchableOpacity
                    onPress={() =>
                      handleEditProperty(prop.propertyId, prop.propertyName)
                    }
                  >
                    <Text style={styles.details}>Edit</Text>
                  </TouchableOpacity>
                  <TouchableOpacity
                    onPress={() => handleDeleteProperty(prop.propertyId)}
                  >
                    <Text style={[styles.details, { color: "red" }]}> 
                      Delete
                    </Text>
                  </TouchableOpacity>
                </View>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>No properties</Text>
          )}

          <TextInput
            style={styles.input}
            placeholder="New Property Name"
            value={newPropertyName}
            onChangeText={setNewPropertyName}
          />
          <TouchableOpacity style={styles.button} onPress={handleAddProperty}>
            <Text style={styles.buttonText}>+ Add Property</Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </LinearGradient>
  );
}