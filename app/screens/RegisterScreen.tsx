import { useRouter } from "expo-router";
import React, { useState } from "react";
import { Alert, Button, StyleSheet, TextInput, View } from "react-native";
import { createUser } from "../../apis/users";

export default function RegisterScreen() {
  const [userId, setUserId] = useState(0);
  const [userFirstName, setFirstName] = useState("");
  const [userLastName, setLastName] = useState("");
  const [userEmailAddress, setEmail] = useState("");
  const [userPhoneNumber, setPhoneNumber] = useState("");
  const [userPassword, setPassword] = useState("");
  const [userRole, setRole] = useState("");
  
  const router = useRouter();

  async function handleRegister() {
    try {
      await createUser({
        userId,
        userFirstName,
        userLastName,
        userEmailAddress,
        userPhoneNumber,
        userPassword,
        userRole
      });
      Alert.alert("Success", "Account created!");
      router.replace("/screens/LoginScreen");
    } catch (err: any) {
      Alert.alert("Error", err.message);
    }
  }

  return (
    <View style={styles.container}>
      <TextInput
        style={styles.input}
        placeholder="First Name"
        placeholderTextColor="#909090ff"
        value={userFirstName}
        onChangeText={setFirstName}
      />
      <TextInput
        style={styles.input}
        placeholder="Last Name"
        placeholderTextColor="#909090ff"
        value={userLastName}
        onChangeText={setLastName}
      />
      <TextInput
        style={styles.input}
        placeholder="Email"
        placeholderTextColor="#909090ff"
        value={userEmailAddress}
        onChangeText={setEmail}
        autoCapitalize="none"
      />
      <TextInput
        style={styles.input}
        placeholder="Phone Number"
        placeholderTextColor="#909090ff"
        value={userPhoneNumber}
        onChangeText={setPhoneNumber}
        keyboardType="phone-pad"
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        placeholderTextColor="#909090ff"
        value={userPassword}
        onChangeText={setPassword}
        secureTextEntry
      />
      <Button title="Register" onPress={handleRegister} />
      <Button title="Back to Login" onPress={() => router.back()} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    justifyContent: "center", 
    padding: 20,
    backgroundColor: "#fff", // <-- add this
  },
input: {
  borderWidth: 1,
  borderColor: "#ccc",
  borderRadius: 8,
  marginBottom: 10,
  padding: 10,
  backgroundColor: "#fff", // ensures input itself is white
  color: "#000",           // ensures text is readable
}

});
