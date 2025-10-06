import { MaterialIcons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useState } from "react";
import { FlatList, StyleSheet, Text, View } from "react-native";
import apiRequest from "../../apis/client";

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
      const overdueData = await apiRequest("/overduePayments");
      setOverdue(overdueData);
      const upcomingData = await apiRequest("/upcomingPayments");
      setUpcoming(upcomingData);
    }
    fetchData();
  }, []);

  const formatDate = (unix?: number) => {
    if (!unix) return "No payment";
    const date = new Date(unix * 1000);
    return date.toLocaleDateString();
  };

  const renderItem = ({ item }: { item: RentPayment }) => {
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
      <View style={styles.row}>
        <View style={styles.rowLeft}>
          <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
          <Text style={styles.address}>{item.address}{item.unit ? ` #${item.unit}` : ""}</Text>
        </View>

        <View style={styles.rowRight}>
          <View style={{ flexDirection: "row", alignItems: "center" }}>
            {icon}
            <Text style={{ ...styles.amount, color, marginRight: 8 }}>
              ${item.rentAmount.toFixed(2)}
            </Text>
           
          </View>
           <Text style={{ color: "#555", marginRight: 8  }}>{formatDate(item.lastPaymentUnix)}</Text>
        </View>
      </View>

    );
  };

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <FlatList
        ListHeaderComponent={
          <>
            <Text style={styles.title}>Overdue Rent</Text>
            <View style={styles.card}>
              <FlatList
                data={overdue}
                renderItem={renderItem}
                keyExtractor={(item, idx) => `overdue-${idx}`}
              />
            </View>
            <Text style={styles.title}>Upcoming Rent</Text>
            <View style={styles.card}>
              <FlatList
                data={upcoming}
                renderItem={renderItem}
                keyExtractor={(item, idx) => `upcoming-${idx}`}
              />
            </View>
          </>
        }
        data={[]}
        renderItem={null}
      />
    </LinearGradient>
  );
}

const styles = StyleSheet.create({
  gradient: {
    flex: 1,
    padding: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: "bold",
    color: "#fff",
    marginBottom: 12,
    marginTop: 50,
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
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    borderBottomWidth: 1,
    borderBottomColor: "#ddd",
    paddingVertical: 8,
  },
  rowLeft: {
    flexDirection: "column",
    flexShrink: 1,
  },
  rowRight: {
    justifyContent: "center",
    alignItems: "flex-end",
  },
  name: {
    fontSize: 16,
    fontWeight: "600",
    color: "#222",
  },
  address: {
    color: "#555",
  },
  amount: {
    fontWeight: "bold",
  },
});
