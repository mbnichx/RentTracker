// app/DashboardScreen.js
import { signOut } from "firebase/auth";
import { Button, Text, View } from "react-native";
import { auth } from "../firebaseConfig";

export default function DashboardScreen() {
  return (
    <View style={styles.container}>
      <Text>Welcome to the Dashboard!</Text>
      <Button title="Logout" onPress={() => signOut(auth)} />
    </View>
  );
}

// const styles = StyleSheet.create({
//   container: { flex: 1, justifyContent: "center", alignItems: "center" },
// });
