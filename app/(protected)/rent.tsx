/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Rent screen — shows overdue and upcoming rent payments. This file defines a
 * small `RentPayment` type used for local rendering, fetches two endpoints
 * (`/overduePayments` and `/upcomingPayments`), and renders each row with a
 * small icon indicating status.
 */
import { MaterialIcons } from "@expo/vector-icons";
import { useFocusEffect } from "@react-navigation/native";
import { LinearGradient } from "expo-linear-gradient";
import React, { useCallback, useEffect, useState } from "react";
import { Button, Modal, ScrollView, Text, TextInput, TouchableOpacity, View } from "react-native";
import apiRequest from "../../apis/client";
import { styles } from "./style";

// Local rendering shape for payments. The backend may return additional
// fields; adapt this type if you add stricter typing across the app.
type RentPayment = {
  firstName: string;
  lastName: string;
  address: string;
  unit: string;
  leaseId?: number;
  rentAmount: number;
  paymentStatus: string; // "Overdue", "Due", "Paid"
  lastPaymentUnix?: number;
};

export default function RentScreen() {
  // Modal state for marking as paid
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedPayment, setSelectedPayment] = useState<RentPayment | null>(null);
  const [notes, setNotes] = useState("");
  const [photo, setPhoto] = useState<string | null>(null);
  const [overdue, setOverdue] = useState<RentPayment[]>([]);
  const [upcoming, setUpcoming] = useState<RentPayment[]>([]);

  // Fetch overdue and upcoming payments. Extracted so it can be
  // re-used after actions (e.g., marking a payment as paid).
  const fetchData = useCallback(async () => {
    try {
      const overdueData = await apiRequest("/overduePayments");
      // Normalize incoming rows to always provide a numeric `leaseId` that
      // the UI and modal can rely on. Backend may return `leaseId`,
      // `lease_id`, or `leaseID` depending on the view; map them here.
      const normalizedOverdue = (overdueData || []).map((x: any) => ({
        ...x,
        leaseId: x.leaseId ?? x.lease_id ?? x.leaseID ?? x.leaseid,
      }));
      setOverdue(normalizedOverdue);
      const upcomingData = await apiRequest("/upcomingPayments");
      const normalizedUpcoming = (upcomingData || []).map((x: any) => ({
        ...x,
        leaseId: x.leaseId ?? x.lease_id ?? x.leaseID ?? x.leaseid,
      }));
      setUpcoming(normalizedUpcoming);
    } catch (err) {
      console.error("Rent fetch error:", err);
    }
  }, []);

  useEffect(() => {
    // Fetch on mount
    fetchData();
  }, [fetchData]);

  // Refresh whenever this screen gains focus (navigate to/from)
  useFocusEffect(
    useCallback(() => {
      fetchData();
    }, [fetchData])
  );

  // Helper to format a Unix timestamp (in seconds) into a readable date
  const formatDate = (unix?: number) => {
    if (!unix) return "No payment";
    const date = new Date(unix * 1000);
    return date.toLocaleDateString();
  };

  // Renders a single row for the payments list. The icon and color vary by
  // payment status so the user can quickly scan the list.
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

    // Only Overdue rows are clickable for marking as paid
    const RowComponent = item.paymentStatus === "Overdue" ? TouchableOpacity : View;
    return (
      <RowComponent
        key={idx}
        style={styles.row}
        onPress={item.paymentStatus === "Overdue" ? () => {
          // Ensure selectedPayment always has leaseId available (map again in case)
          setSelectedPayment({ ...item, leaseId: item.leaseId });
          setModalVisible(true);
        } : undefined}
      >
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
      </RowComponent>
    );
  };

  // Handle payment submission
  const handleMarkPaid = async () => {
    if (!selectedPayment) return;
    try {
      // Ensure we have a leaseId to associate the payment with
      const leaseId = selectedPayment.leaseId;
 
      if (!leaseId) {
        console.error("Selected payment item does not include a leaseId");
        return;
      }

      // Compose Payment object for backend. Send photo as base64 string; the
      // backend will receive this as a base64-encoded field and can decode it.
      const paymentPayload = {
        leaseId: leaseId,
        paymentAmount: selectedPayment.rentAmount,
        paymentDateUnix: Math.floor(Date.now() / 1000),
        paymentNotes: notes,
        // paymentConfirmation: photo || undefined,
      };
      console.log(paymentPayload)
      await apiRequest("/payments", "POST", paymentPayload);
      setModalVisible(false);
      setNotes("");
      // setPhoto(null);
      setSelectedPayment(null);
      // Refresh lists so UI reflects the newly recorded payment
      try {
        await fetchData();
      } catch (e) {
        // swallow; fetchData already logs errors but keep safety here
        console.error("Refresh after payment failed:", e);
      }
    } catch (err) {
      console.error("Payment save error:", err);
    }
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

      {/* Modal for marking as paid */}
      <Modal
        visible={modalVisible}
        animationType="slide"
        transparent={true}
        onRequestClose={() => setModalVisible(false)}
      >
        <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#0008' }}>
          <View style={{ backgroundColor: '#fff', borderRadius: 16, padding: 24, width: '90%' }}>
            <Text style={{ fontSize: 20, fontWeight: 'bold', marginBottom: 12 }}>Mark as Paid</Text>
            <Text style={{ marginBottom: 8 }}>Amount: ${selectedPayment?.rentAmount.toFixed(2)}</Text>
            <TextInput
              placeholder="Notes (optional)"
              value={notes}
              onChangeText={setNotes}
              style={{ borderWidth: 1, borderColor: '#ccc', borderRadius: 8, padding: 8, marginBottom: 12 }}
            />
            <View style={{ flexDirection: 'row', justifyContent: 'space-between', marginTop: 18 }}>
              <Button title="Cancel" color="#888" onPress={() => setModalVisible(false)} />
              <Button title="Mark as Paid" color="#28a745" onPress={handleMarkPaid} />
            </View>
          </View>
        </View>
      </Modal>
    </LinearGradient>
  );
}
