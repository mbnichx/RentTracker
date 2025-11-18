/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Tenants screen — displays a list of current tenants. The screen fetches
 * lease overview data from the backend which includes tenant name and their
 * unit information. The component keeps the UI intentionally simple and
 * renders the returned objects directly; add stronger typing if you
 * formalize the API shapes later.
 */
import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View } from "react-native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function TenantsScreen() {
  // Tenant array returned by the /leaseOverview endpoint
  const [tenants, setTenants] = useState<any[]>([]);

  useEffect(() => {
    const fetchTenants = async () => {
      try {
        // Lease overview endpoint contains tenant contact and lease dates
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
