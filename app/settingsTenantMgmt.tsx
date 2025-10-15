import { useRoute } from "@react-navigation/native";
import { LinearGradient } from "expo-linear-gradient";
import { router } from "expo-router";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, TextInput, TouchableOpacity, View } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import apiRequest from "../apis/client";
import { styles } from "./(protected)/style";

import * as u from "../apis/units";

export default function TenantManagementScreen() {
  const route = useRoute();
  const { propertyId } = route.params as { propertyId: number };
  const [units, setUnits] = useState<any[]>([]);
  const [openOccupancy, setOpenOccupancy] = useState(false);
  const [occupancyValue, setOccupancyValue] = useState(null);

  const [openPets, setOpenPets] = useState(false);
  const [petsValue, setPetsValue] = useState(null);
  const [petsItems, setPetsItems] = useState([
    { label: "No Pets", value: "No Pets" },
    { label: "Cats", value: "Cats" },
    { label: "Dogs", value: "Dogs" },
    { label: "Other", value: "Other" },
  ]);

  useEffect(() => {
    fetchUnits();
  }, [propertyId]);

  const fetchUnits = async () => {
    try {
      const res = await u.getUnits(propertyId);
      setUnits(res || []);
    } catch (err) {
      console.error("Units fetch error:", err);
    }
  };

  const updateTenantInfo = async (unitId: string, field: string, value: any) => {
    try {
      const unit = units.find((u) => u.id === unitId);
      const updatedTenant = { ...unit.tenant, [field]: value };
      await apiRequest(`/tenants`, "PUT", updatedTenant);
      setUnits((prev) =>
        prev.map((u) => (u.id === unitId ? { ...u, tenant: updatedTenant } : u))
      );
    } catch (err) {
      console.error("Tenant update error:", err);
    }
  };

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <View style={{ flex: 1, padding: 20 }}>
          <TouchableOpacity onPress={() => router.back()}>
            <Text style={{ color: "#2575fc", fontSize: 16 }}>← Back to Settings</Text>
          </TouchableOpacity>
        </View>
        <Text style={styles.title}>Tenants</Text>
        {units.length > 0 ? (
          units.map((unit, idx) => (
            <View key={idx} style={styles.card}>
              <Text style={styles.cardTitle}>Unit {unit.propertyUnitNumber}</Text>

              <TextInput
                style={styles.input}
                placeholder="Tenant Name"
                placeholderTextColor="#999"
                value={unit.tenant?.name || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "name", val)}
              />
              <TextInput
                style={styles.input}
                placeholder="Tenant Email"
                placeholderTextColor="#999"
                value={unit.tenant?.email || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "email", val)}
              />
              <TextInput
                style={styles.input}
                placeholder="Tenant Phone"
                placeholderTextColor="#999"
                value={unit.tenant?.phone || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "phone", val)}
              />

              <DropDownPicker
                open={openOccupancy}
                value={occupancyValue}
                items={[
                  { label: "Vacant", value: "Vacant" },
                  { label: "Occupied", value: "Occupied" },
                  { label: "Pending", value: "Pending" },
                ]}
                setOpen={setOpenOccupancy}
                setValue={(cb) =>
                  updateTenantInfo(unit.id, "occupancy", cb(occupancyValue))
                }
                placeholder="Select Occupancy"

                style={styles.dropdown}
                dropDownContainerStyle={styles.dropdownContainer}
                zIndex={2000} // must be lower than occupancy to avoid overlap issues

              />
              <DropDownPicker
                open={openPets}
                value={petsValue}
                items={petsItems}
                setOpen={setOpenPets}
                setValue={setPetsValue}
                setItems={setPetsItems}
                placeholder="Select Pets"
                style={styles.dropdown}
                dropDownContainerStyle={styles.dropdownContainer}
                zIndex={1000} // must be lower than occupancy to avoid overlap issues
              />

            </View>
          ))
        ) : (
          <Text style={styles.emptyText}>No units found for this property</Text>
        )}
      </ScrollView>
    </LinearGradient>
  );
}
