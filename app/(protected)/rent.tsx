import { MaterialIcons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { ScrollView, Text, View } from "react-native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

type RentPayment = {
  firstName: string;
  lastName: string;
  address: string;
  unit: string;
  rentAmount: number;
  paymentStatus: string; // "Overdue", "Due", "Paid"
  lastPaymentUnix?: number;
};

export default function RentScreen() {
  const [overdue, setOverdue] = useState<RentPayment[]>([]);
  const [upcoming, setUpcoming] = useState<RentPayment[]>([]);

  useEffect(() => {
    async function fetchData() {
      try {
        const overdueData = await apiRequest("/overduePayments");
        setOverdue(overdueData || []);
        const upcomingData = await apiRequest("/upcomingPayments");
        setUpcoming(upcomingData || []);
      } catch (err) {
        console.error("Rent fetch error:", err);
      }
    }
    fetchData();
  }, []);

  const formatDate = (unix?: number) => {
    if (!unix) return "No payment";
    const date = new Date(unix * 1000);
    return date.toLocaleDateString();
  };

  const renderRow = (item: RentPayment, idx: number) => {
    let color = "#000";
    let icon = null;

    if (item.paymentStatus === "Overdue") {
      color = "#d32f2f";
      icon = <MaterialIcons name="warning" size={20} color={color} style={{ marginRight: 5 }} />;
    } else if (item.paymentStatus === "Paid") {
      color = "#28a745";
      icon = <MaterialIcons name="check-circle" size={20} color={color} style={{ marginRight: 5 }} />;
    }

    return (
      <View key={idx} style={styles.row}>
        {/* Left side: name + address */}
        <View style={styles.rowLeft}>
          <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
          <Text style={styles.address}>{item.address}{item.unit ? ` #${item.unit}` : ""}</Text>
        </View>

        {/* Right side: amount + date */}
        <View style={styles.rowRight}>
          <View style={{ flexDirection: "row", alignItems: "center" }}>
            {icon}
            <Text style={{ ...styles.amount, color, marginRight: 8 }}>
              ${item.rentAmount.toFixed(2)}
            </Text>
          </View>
          <Text style={{ color: "#555" }}>{formatDate(item.lastPaymentUnix)}</Text>
        </View>
      </View>
    );
  };

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <Text style={styles.title}>Rent</Text>

        {/* Overdue Section */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Overdue Rent</Text>
          {overdue.length > 0 ? (
            overdue.map(renderRow)
          ) : (
            <Text style={styles.emptyText}>No overdue rent</Text>
          )}
        </View>

        {/* Upcoming Section */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Upcoming Rent</Text>
          {upcoming.length > 0 ? (
            upcoming.map(renderRow)
          ) : (
            <Text style={styles.emptyText}>No upcoming rent</Text>
          )}
        </View>
      </ScrollView>
    </LinearGradient>
  );
}
