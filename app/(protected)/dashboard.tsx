import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ActivityIndicator, ScrollView, StyleSheet, Text, View } from "react-native";
import apiRequest from "../../apis/client"; // 👈 your centralized API request helper

export default function DashboardScreen() {
  const [overduePayments, setOverduePayments] = useState<any[]>([]);
  const [maintenanceRequests, setMaintenanceRequests] = useState<any[]>([]);
  const [leases, setLeases] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
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
        console.error("Dashboard fetch error:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  if (loading) {
    return (
      <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
        <View style={styles.centered}>
          <ActivityIndicator size="large" color="#fff" />
          <Text style={styles.loadingText}>Loading dashboard...</Text>
        </View>
      </LinearGradient>
    );
  }

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
                  <Text style={styles.amount}>${item.rentAmount}</Text>
                </View>
              </View>

            ))
          ) : (
            <Text style={styles.emptyText}>No overdue payments 🎉</Text>
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
            <Text style={styles.emptyText}>No maintenance requests 🧰</Text>
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

const styles = StyleSheet.create({
  gradient: {
    flex: 1,
  },
  scrollContainer: {
    padding: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: "bold",
    color: "#fff",
    marginBottom: 24,
    marginTop: 50,
    textAlign: "center",
  },
  card: {
    backgroundColor: "rgba(255, 255, 255, 0.9)",
    borderRadius: 16,
    padding: 16,
    marginBottom: 20,
    shadowColor: "#000",
    shadowOpacity: 0.2,
    shadowRadius: 6,
    elevation: 3,
  },
  cardTitle: {
    fontSize: 20,
    fontWeight: "600",
    marginBottom: 12,
    color: "#333",
  },
  name: {
    fontSize: 16,
    fontWeight: "600",
    color: "#222",
  },
  details: {
    color: "#555",
  },
  amount: {
    color: "#d32f2f",
    fontWeight: "bold",
  },
  address: {
    fontSize: 14,
    color: "#555",
    marginTop: 2,
  },
  emptyText: {
    fontStyle: "italic",
    color: "#555",
  },
  row: {
    flexDirection: "row",      // horizontal layout
    justifyContent: "space-between", // left vs right
    alignItems: "center",      // vertically center right content
    borderBottomWidth: 1,
    borderBottomColor: "#ddd",
    paddingVertical: 8,
  },
  rowLeft: {
    flexDirection: "column",   // stack name and address vertically
    flexShrink: 1,             // allows long addresses to wrap
  },
  rowRight: {
    justifyContent: "center",  // vertically center
    alignItems: "flex-end",    // right align
  },
  centered: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  loadingText: {
    color: "#fff",
    marginTop: 10,
  },
});
