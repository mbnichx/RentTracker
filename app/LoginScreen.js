// app/LoginScreen.js
import { signInWithEmailAndPassword } from "firebase/auth";
import { useState } from "react";
import { Alert, Button, View } from "react-native";
import FormInput from "../components/FormInput";
import { auth } from "../firebaseConfig";

export default function LoginScreen({ navigation }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    try {
      await signInWithEmailAndPassword(auth, email, password);
      // onAuthStateChanged in App.js handles redirect
    } catch (error) {
      Alert.alert("Login failed", error.message);
    }
  };

  return (
    <View style={styles.container}>
      <FormInput placeholder="Email" value={email} onChangeText={setEmail} />
      <FormInput placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry />
      <Button title="Login" onPress={handleLogin} />
      <Button title="Go to Sign Up" onPress={() => navigation.navigate("Signup")} />
    </View>
  );
}

// const styles = StyleSheet.create({
//   container: { flex: 1, padding: 20, justifyContent: "center" },
// });
