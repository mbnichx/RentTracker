import { StyleSheet } from "react-native";

export const styles = StyleSheet.create({
  gradient: {
    flex: 1,
    padding: 20,
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
  cardTitle: {
    fontSize: 20,
    fontWeight: "600",
    marginBottom: 12,
    color: "#333",
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
  centered: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
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
  overdueAmount: {
    fontWeight: "bold",
    color: "#d32f2f",
  },
  details: {                  // Relevant info., ex. maintenance request description
    color: "#555",
  },
  emptyText: {                // Placeholder text
    fontStyle: "italic",
    color: "#555",
  },
  filterHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 10,
  },
  dropdownCompact: {
    backgroundColor: "rgba(254, 254, 254, 1)",
    borderWidth: 0,
    borderRadius: 8,
    minHeight: 35,
    width: 130,
  },
  dropdownTextCompact: {
    color: "#171515ff",
    fontSize: 13,
  },
  dropdownContainerCompact: {
    backgroundColor: "rgba(254, 254, 254, 1)",
    borderWidth: 0,
    borderRadius: 8,
  },
  status: {
    marginTop: 6,
    backgroundColor: "rgba(255, 255, 255, 0.2)",
    color: "#fff",
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 6,
    fontSize: 12,
    fontWeight: "600",
    alignSelf: "flex-end",
    overflow: "hidden",
  },
  input: {
    backgroundColor: "#fff",
    color: "#000000",
    borderRadius: 8,
    padding: 10,
    marginBottom: 10,
  },
  button: {
    backgroundColor: "#2575fc",
    padding: 12,
    borderRadius: 8,
    alignItems: "center",
    marginTop: 5,
  },
  buttonText: {
    color: "#fff",
    fontWeight: "bold",
  },
  picker: {
    backgroundColor: "#fff",
    color: "#999",
    borderRadius: 8,
    marginBottom: 10,
  },
  dropdown: {
    backgroundColor: "#fff",
    borderColor: "#ccc",
    borderRadius: 8,
    marginBottom: 10,
  },
  dropdownContainer: {
    backgroundColor: "#fff",
    borderColor: "#ccc",
  },

});
