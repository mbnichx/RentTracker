/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Dashboard screen — aggregates several small datasets used to give the user a
 * quick overview (overdue payments, maintenance requests, and active leases).
 *
 * This component uses `apiRequest` directly to fetch three endpoints in
 * parallel. The responses are stored in local state and rendered into simple
 * list cards. The screen is intentionally lightweight — the individual items
 * are rendered inline and assume the backend returns objects with the
 * fields referenced below (firstName, lastName, address, unit, etc.).
 */
import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View } from "react-native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function DashboardScreen() {
  // Local UI state for each dataset. Using `any[]` keeps this file simple
  // — consider adding typed response shapes for stricter type-safety later.
  const [overduePayments, setOverduePayments] = useState<any[]>([]);
  const [maintenanceRequests, setMaintenanceRequests] = useState<any[]>([]);
  const [leases, setLeases] = useState<any[]>([]);

  useEffect(() => {
    // Fetch the three dashboard endpoints in parallel for speed.
    const fetchDashboardData = async () => {
      try {
        const [overdueRes, maintenanceRes, leasesRes] = await Promise.all([
          apiRequest("/overduePayments", "GET"),
          apiRequest("/maintenanceRequestStatus", "GET"),
          apiRequest("/leaseOverview", "GET"),
        ]);

        setOverduePayments(overdueRes || []);
        setMaintenanceRequests(maintenanceRes || []);
        setLeases(leasesRes || []);
      } catch (err) {
        // Keep the UI stable and log the issue for debugging
        console.error("Dashboard fetch error:", err);
      } 
    };

    fetchDashboardData();
  }, []);

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <Text style={styles.title}>Dashboard</Text>

        {/* Overdue Payments */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Overdue Rent Payments</Text>
          {overduePayments.length > 0 ? (
            overduePayments.map((item, idx) => (
              <View key={idx} style={styles.row}>
                {/* Left side: name + address stacked */}
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
                  <Text style={styles.address}>{item.address}{item.unit ? ` #${item.unit}` : ""}</Text>
                </View>

                {/* Right side: rent amount / date, vertically centered */}
                <View style={styles.rowRight}>
                  <Text style={styles.overdueAmount}>${item.rentAmount.toFixed(2)}</Text>
                </View>
              </View>

            ))
          ) : (
            <Text style={styles.emptyText}>No overdue payments</Text>
          )}

        </View>

        {/* Maintenance Requests */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Maintenance Requests</Text>
          {maintenanceRequests.length > 0 ? (
            maintenanceRequests.map((item, idx) => (
              <View key={idx} style={styles.row}>
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
                  <Text style={styles.address}>{item.address}{item.unit ? ` #${item.unit}` : ""}</Text>
                </View>
                <View style={styles.rowRight}>
                  <Text style={styles.details}>{new Date(item.dateCreated).toLocaleDateString()}</Text>
                  <Text style={styles.details}>{item.description}</Text>
                </View>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>No maintenance requests</Text>
          )}
        </View>

        {/* Leases */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Active Leases</Text>
          {leases.length > 0 ? (
            leases.map((item, idx) => (
              <View key={idx} style={styles.row}>
                <View style={styles.rowLeft}>
                  <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
                  <Text style={styles.address}>{item.address}{item.unit ? ` #${item.unit}` : ""}</Text>
                </View>
                <View style={styles.rowRight}>
                  <Text style={styles.details}>{new Date(item.leaseStartDate).toLocaleDateString()}</Text>
                </View>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>No active leases 🏠</Text>
          )}
        </View>
      </ScrollView>
    </LinearGradient>
  );
}
