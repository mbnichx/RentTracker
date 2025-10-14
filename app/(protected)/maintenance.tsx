import { LinearGradient } from "expo-linear-gradient";
import React, { useEffect, useMemo, useState } from "react";
import { ScrollView, Text, View } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function MaintenanceScreen() {
  const [currMaintenanceReqs, setCurrMaintenanceReqs] = useState<any[]>([]);
  const [selectedStatus, setSelectedStatus] = useState<string>("All");

  // dropdown state
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState([
    { label: "All", value: "All" },
    { label: "Open", value: "Open" },
    { label: "In Progress", value: "In Progress" },
    { label: "Closed", value: "Closed" },
  ]);

  useEffect(() => {
    const fetchMaintenanceData = async () => {
      try {
        const [currMaintenanceReqs] = await Promise.all([
          apiRequest("/maintenanceRequestStatus", "GET"),
        ]);
        setCurrMaintenanceReqs(currMaintenanceReqs || []);
      } catch (err) {
        console.error("Maintenance fetch error:", err);
      }
    };

    fetchMaintenanceData();
  }, []);

  const activeRequests = useMemo(() => {
    return currMaintenanceReqs.filter(
      (item) =>
        item.status?.toLowerCase() === "open" ||
        item.status?.toLowerCase() === "in progress"
    );
  }, [currMaintenanceReqs]);

  const filteredRequests = useMemo(() => {
    if (selectedStatus === "All") return currMaintenanceReqs;
    return currMaintenanceReqs.filter(
      (item) =>
        item.status?.toLowerCase() === selectedStatus.toLowerCase()
    );
  }, [currMaintenanceReqs, selectedStatus]);

  return (
    <LinearGradient colors={["#6a11cb", "#2575fc"]} style={styles.gradient}>
      <View style={styles.scrollContainer}>
        <Text style={styles.title}>Maintenance</Text>
        <ScrollView
          contentContainerStyle={{ paddingBottom: 40 }}
          showsVerticalScrollIndicator={false}
        >
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Active Requests</Text>
            {activeRequests.length > 0 ? (
              activeRequests.map((item, idx) => (
                <View key={idx} style={styles.row}>
                  <View style={styles.rowLeft}>
                    <Text style={styles.name}>
                      {item.firstName} {item.lastName}
                    </Text>
                    <Text style={styles.address}>
                      {item.address}
                      {item.unit ? ` #${item.unit}` : ""}
                    </Text>
                  </View>
                  <View style={styles.rowRight}>
                    <Text style={styles.details}>
                      {new Date(item.dateCreated).toLocaleDateString()}
                    </Text>
                    <Text style={styles.details}>{item.description}</Text>
                    <Text
                      style={[
                        styles.status,
                        item.status === "Open" && { backgroundColor: "rgba(76, 175, 80, 0.4)" },
                        item.status === "In Progress" && { backgroundColor: "rgba(255, 193, 7, 0.4)" },
                        item.status === "Closed" && { backgroundColor: "rgba(244, 67, 54, 0.4)" },
                      ]}
                    >
                      {item.status}
                    </Text>
                  </View>
                </View>
              ))
            ) : (
              <Text style={styles.emptyText}>No active requests</Text>
            )}
          </View>

          <View style={styles.card}>
            <View style={styles.filterHeader}>
              <Text style={styles.cardTitle}>All Requests</Text>
            </View>
            <View>
              <DropDownPicker
                open={open}
                value={selectedStatus}
                items={items}
                setOpen={setOpen}
                setValue={setSelectedStatus}
                setItems={setItems}
                style={styles.dropdownCompact}
                textStyle={styles.dropdownTextCompact}
                dropDownContainerStyle={styles.dropdownContainerCompact}
                placeholder="Filter"
                listMode="SCROLLVIEW"
                zIndex={1000}
              />
            </View>

            {filteredRequests.length > 0 ? (
              filteredRequests.map((item, idx) => (
                <View key={idx} style={styles.row}>
                  <View style={styles.rowLeft}>
                    <Text style={styles.name}>
                      {item.firstName} {item.lastName}
                    </Text>
                    <Text style={styles.address}>
                      {item.address}
                      {item.unit ? ` #${item.unit}` : ""}
                    </Text>
                  </View>
                  <View style={styles.rowRight}>
                    <Text style={styles.details}>
                      {new Date(item.dateCreated).toLocaleDateString()}
                    </Text>
                    <Text style={styles.details}>{item.description}</Text>
                    <Text
                      style={[
                        styles.status,
                        item.status === "Open" && { backgroundColor: "rgba(76, 175, 80, 0.4)" },
                        item.status === "In Progress" && { backgroundColor: "rgba(255, 193, 7, 0.4)" },
                        item.status === "Closed" && { backgroundColor: "rgba(244, 67, 54, 0.4)" },
                      ]}
                    >
                      {item.status}
                    </Text>
                  </View>
                </View>
              ))
            ) : (
              <Text style={styles.emptyText}>No requests found</Text>
            )}
          </View>
        </ScrollView>
      </View>
    </LinearGradient>
  );
}
