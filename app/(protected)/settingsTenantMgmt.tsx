import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View, TextInput } from "react-native";
import { Picker } from "@react-native-picker/picker";
import { useRoute } from "@react-navigation/native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function TenantManagementScreen() {
  const route = useRoute();
  const { propertyId } = route.params as { propertyId: string };
  const [units, setUnits] = useState<any[]>([]);

  useEffect(() => {
    fetchUnits();
  }, [propertyId]);

  const fetchUnits = async () => {
    try {
      const res = await apiRequest(`/properties/${propertyId}/units`, "GET");
      setUnits(res || []);
    } catch (err) {
      console.error("Units fetch error:", err);
    }
  };

  const updateTenantInfo = async (unitId: string, field: string, value: any) => {
    try {
      const unit = units.find((u) => u.id === unitId);
      const updatedTenant = { ...unit.tenant, [field]: value };
      await apiRequest(`/units/${unitId}/tenant`, "PUT", updatedTenant);
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
        <Text style={styles.title}>Tenants</Text>
        {units.length > 0 ? (
          units.map((unit, idx) => (
            <View key={idx} style={styles.card}>
              <Text style={styles.cardTitle}>Unit {unit.unitNumber}</Text>

              <TextInput
                style={styles.input}
                placeholder="Tenant Name"
                value={unit.tenant?.name || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "name", val)}
              />
              <TextInput
                style={styles.input}
                placeholder="Tenant Email"
                value={unit.tenant?.email || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "email", val)}
              />
              <TextInput
                style={styles.input}
                placeholder="Tenant Phone"
                value={unit.tenant?.phone || ""}
                onChangeText={(val) => updateTenantInfo(unit.id, "phone", val)}
              />

              <Picker
                selectedValue={unit.tenant?.occupancy || "Vacant"}
                onValueChange={(val) => updateTenantInfo(unit.id, "occupancy", val)}
                style={styles.picker}
              >
                <Picker.Item label="Vacant" value="Vacant" />
                <Picker.Item label="Occupied" value="Occupied" />
                <Picker.Item label="Pending" value="Pending" />
              </Picker>

              <Picker
                selectedValue={unit.tenant?.pets || "No Pets"}
                onValueChange={(val) => updateTenantInfo(unit.id, "pets", val)}
                style={styles.picker}
              >
                <Picker.Item label="No Pets" value="No Pets" />
                <Picker.Item label="Cats" value="Cats" />
                <Picker.Item label="Dogs" value="Dogs" />
                <Picker.Item label="Other" value="Other" />
              </Picker>
            </View>
          ))
        ) : (
          <Text style={styles.emptyText}>No units found for this property</Text>
        )}
      </ScrollView>
    </LinearGradient>
  );
}
