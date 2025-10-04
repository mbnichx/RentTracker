import { Image } from 'expo-image';
import { Platform, StyleSheet } from 'react-native';
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import React, { useEffect, useState } from "react";
import { View, Text, FlatList, Button, ActivityIndicator } from "react-native";
import { Collapsible } from '@/components/ui/collapsible';
import { ExternalLink } from '@/components/external-link';
import ParallaxScrollView from '@/components/parallax-scroll-view';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Fonts } from '@/constants/theme';

import { getTenants } from "../apis/tenants";


type RootStackParamList = {
  Tenants: undefined;
  AddTenant: undefined;
};

type AddTenantScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "AddTenant"
>;

type Props = {
  navigation: AddTenantScreenNavigationProp;
};

export default function TenantsScreen({ navigation }: Props) {
  const [tenants, setTenants] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTenants = async () => {
      try {
        const data = await getTenants();
        setTenants(data);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchTenants();
  }, []);

  if (loading) return <ActivityIndicator size="large" />;

  return (
    <View style={{ flex: 1, padding: 20 }}>
      <Button title="Add Tenant" onPress={() => navigation.navigate("AddTenant")} />
      <FlatList
        data={tenants}
        keyExtractor={(item) => item.tenantId.toString()}
        renderItem={({ item }) => (
          <View style={{ padding: 10, borderBottomWidth: 1 }}>
            <Text>{item.tenantFirstName} {item.tenantLastName}</Text>
            <Text>{item.tenantEmailAddress}</Text>
            <Text>{item.tenantPhoneNumber}</Text>
          </View>
        )}
      />
    </View>
  );
}