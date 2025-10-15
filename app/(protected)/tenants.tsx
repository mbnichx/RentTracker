import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View } from "react-native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function TenantsScreen() {
  const [tenants, setTenants] = useState<any[]>([]);

  useEffect(() => {
    const fetchTenants = async () => {
      try {
        const res = await apiRequest("/leaseOverview", "GET");
        setTenants(res || []);
      } catch (err) {
        console.error("Tenants fetch error:", err);
      }
    };

    fetchTenants();
  }, []);

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <Text style={styles.title}>Tenants</Text>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>Current Tenants</Text>
          {tenants.length > 0 ? (
            tenants.map((tenant, idx) => (
              <View key={idx} style={styles.row}>
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>
                    {tenant.firstName} {tenant.lastName}
                  </Text>
                  <Text style={styles.address}>
                    {tenant.address}
                    {tenant.unit ? ` #${tenant.unit}` : ""}
                  </Text>
                </View>

                <View style={styles.rowRight}>
                  <Text style={styles.details}>{new Date(tenant.leaseStartDate).toLocaleDateString()}</Text>
                </View>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>No tenants found</Text>
          )}
        </View>
      </ScrollView>
    </LinearGradient>
  );
}
