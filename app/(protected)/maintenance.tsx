/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Maintenance screen — shows active and all maintenance requests with a
 * simple status filter. The UI uses a compact DropDownPicker instance to
 * select the status (All/Open/In Progress/Closed) and memoized lists to
 * compute active and filtered requests.
 */
import { useFocusEffect } from "@react-navigation/native";
import { LinearGradient } from "expo-linear-gradient";
import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Button, Modal, ScrollView, Text, TextInput, TouchableOpacity, View } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import apiRequest from "../../apis/client";
import { styles } from "./style";

export default function MaintenanceScreen() {
  // Full set of maintenance requests returned by the server
  const [currMaintenanceReqs, setCurrMaintenanceReqs] = useState<any[]>([]);
  // Current filter selection shown in the dropdown
  const [selectedStatus, setSelectedStatus] = useState<string>("All");

  // DropDownPicker local state (required by the component)
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState([
    { label: "All", value: "All" },
    { label: "Open", value: "Open" },
    { label: "In Progress", value: "In Progress" },
    { label: "Closed", value: "Closed" },
  ]);

  // Modal state for editing a maintenance request
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedRequestId, setSelectedRequestId] = useState<number | null>(null);
  const [selectedRequestFull, setSelectedRequestFull] = useState<any | null>(null);
  const [modalDescription, setModalDescription] = useState("");
  const [modalStatusOpen, setModalStatusOpen] = useState(false);
  const [modalStatusValue, setModalStatusValue] = useState<string | null>(null);
  const [modalStatusItems, setModalStatusItems] = useState([
    { label: "Open", value: "open" },
    { label: "In Progress", value: "in progress" },
    { label: "Completed", value: "completed" },
  ]);

  const fetchMaintenanceData = useCallback(async () => {
    try {
      // Single endpoint that returns the current status list
      const [currMaintenanceReqs] = await Promise.all([
        apiRequest("/maintenanceRequestStatus", "GET"),
      ]);
      setCurrMaintenanceReqs(currMaintenanceReqs || []);
    } catch (err) {
      console.error("Maintenance fetch error:", err);
    }
  }, []);

  useEffect(() => {
    fetchMaintenanceData();
  }, [fetchMaintenanceData]);

  useFocusEffect(
    useCallback(() => {
      fetchMaintenanceData();
    }, [fetchMaintenanceData])
  );

  // Open edit modal for a specific maintenance request (uses the view's limited data to fetch full record)
  const openEditModal = async (maintenanceRequestId: number) => {
    try {
      // fetch full record from backend
      const full = await apiRequest(`/maintenance/${maintenanceRequestId}`, "GET");
      setSelectedRequestId(maintenanceRequestId);
      setSelectedRequestFull(full);
      setModalDescription(full.maintenanceRequestInfo || "");
      // map backend status values to modal dropdown values
      setModalStatusValue(full.maintenanceRequestStatus || "open");
      setModalVisible(true);
    } catch (err) {
      console.error("Failed to load maintenance request:", err);
    }
  };

  const saveModalEdits = async () => {
    if (!selectedRequestFull) return;
    try {
      const payload = {
        ...selectedRequestFull,
        maintenanceRequestInfo: modalDescription,
        maintenanceRequestStatus: modalStatusValue,
      };
      await apiRequest("/maintenance/update", "PUT", payload);
      setModalVisible(false);
      setSelectedRequestFull(null);
      setSelectedRequestId(null);
      // refresh list
      await fetchMaintenanceData();
    } catch (err) {
      console.error("Failed to save maintenance request:", err);
    }
  };

  // Normalize status for display and color
  const statusInfo = (raw?: string) => {
    const s = (raw || "").toString().toLowerCase().trim();
    if (s === "open") return { label: "Open", bg: "rgba(76, 175, 80, 0.4)" };
    if (s === "in progress") return { label: "In Progress", bg: "rgba(255, 193, 7, 0.4)" };
    if (s === "completed" || s === "closed") return { label: "Completed", bg: "rgba(244, 67, 54, 0.4)" };
    return { label: raw || "", bg: "transparent" };
  };

  // Compute active requests (open + in progress) using useMemo for perf
  const activeRequests = useMemo(() => {
    return currMaintenanceReqs.filter(
      (item) =>
        item.status?.toLowerCase() === "open" ||
        item.status?.toLowerCase() === "in progress"
    );
  }, [currMaintenanceReqs]);

  // Apply the selected status filter or return all when "All" is selected
  const filteredRequests = useMemo(() => {
    if (selectedStatus === "All") return currMaintenanceReqs;
    return currMaintenanceReqs.filter(
      (item) =>
        item.status?.toLowerCase() === selectedStatus.toLowerCase()
    );
  }, [currMaintenanceReqs, selectedStatus]);

  return (
    <>
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
              activeRequests.map((item, idx) => {
                const si = statusInfo(item.status);
                return (
                  <TouchableOpacity key={idx} style={styles.row} onPress={() => openEditModal(item.maintenanceRequestId || item.id || item.requestId)}>
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
                      <Text style={[styles.status, { backgroundColor: si.bg }]}>{si.label}</Text>
                    </View>
                  </TouchableOpacity>
                );
              })
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
              filteredRequests.map((item, idx) => {
                const si = statusInfo(item.status);
                return (
                  <TouchableOpacity key={idx} style={styles.row} onPress={() => openEditModal(item.maintenanceRequestId || item.id || item.requestId)}>
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
                      <Text style={[styles.status, { backgroundColor: si.bg }]}>{si.label}</Text>
                    </View>
                  </TouchableOpacity>
                );
              })
            ) : (
              <Text style={styles.emptyText}>No requests found</Text>
            )}
          </View>
        </ScrollView>
      </View>
    </LinearGradient>
    
    {/* Edit modal */}
    <Modal visible={modalVisible} animationType="slide" transparent={true} onRequestClose={() => setModalVisible(false)}>
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#0008' }}>
        <View style={{ backgroundColor: '#fff', borderRadius: 12, padding: 20, width: '92%' }}>
          <Text style={{ fontSize: 18, fontWeight: 'bold', marginBottom: 8 }}>Edit Request</Text>
          <Text style={{ marginBottom: 6 }}>Description</Text>
          <TextInput value={modalDescription} onChangeText={setModalDescription} multiline style={{ borderWidth: 1, borderColor: '#ccc', borderRadius: 8, padding: 8, minHeight: 80 }} />
          <Text style={{ marginTop: 12, marginBottom: 6 }}>Status</Text>
          <DropDownPicker
            open={modalStatusOpen}
            value={modalStatusValue}
            items={modalStatusItems}
            setOpen={setModalStatusOpen}
            setValue={setModalStatusValue}
            setItems={setModalStatusItems}
            zIndex={2000}
          />
          <View style={{ flexDirection: 'row', justifyContent: 'space-between', marginTop: 16 }}>
            <Button title="Cancel" color="#888" onPress={() => setModalVisible(false)} />
            <Button title="Save" onPress={saveModalEdits} />
          </View>
        </View>
      </View>
    </Modal>
    </>
  );
}
